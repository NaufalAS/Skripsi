import cv2
import numpy as np
from ultralytics import YOLO
import time
import math

# Batas kecepatan (km/h)
limit = 80

# Skala Piksel ke Meter (contoh: 1 meter = 10 piksel)
scale_factor = 0.02  # 1 piksel = 0.02 meter (skala yang disesuaikan)

# Inisialisasi video capture
cap = cv2.VideoCapture("20km.mp4")

# Inisialisasi model YOLO
model = YOLO("yolov8n.pt")  # Model YOLO pre-trained

# Inisialisasi tracker
class EuclideanDistTracker:
    def __init__(self):
        self.center_points = {}  # Menyimpan titik pusat objek
        self.id_count = 0
        self.previous_positions = {}  # Menyimpan posisi kendaraan di frame sebelumnya
        self.s = np.zeros((1, 1000))  # Kecepatan kendaraan
        self.f = np.zeros(1000)  # Flag untuk deteksi objek
        self.count = 0
        self.exceeded = 0  # Hitung jumlah kendaraan yang melebihi batas kecepatan

    def update(self, objects_rect):
        objects_bbs_ids = []
        for rect in objects_rect:
            x, y, w, h = rect
            cx = (x + x + w) // 2
            cy = (y + y + h) // 2

            same_object_detected = False

            for id, pt in self.center_points.items():
                dist = np.hypot(cx - pt[0], cy - pt[1])  # Hitung jarak antar titik

                if dist < 70:  # Jika jarak cukup dekat, anggap objek yang sama
                    self.center_points[id] = (cx, cy)
                    objects_bbs_ids.append([x, y, w, h, id])
                    same_object_detected = True

            # Deteksi objek baru jika belum terdeteksi sebelumnya
            if same_object_detected is False:
                self.center_points[self.id_count] = (cx, cy)
                objects_bbs_ids.append([x, y, w, h, self.id_count])
                self.id_count += 1

        return objects_bbs_ids

    def getsp(self, id, current_position):
        """Menghitung kecepatan kendaraan berdasarkan perbedaan posisi di dua frame berturut-turut"""
        if id in self.previous_positions:
            prev_x, prev_y = self.previous_positions[id]
            current_x, current_y = current_position

            # Hitung jarak Euclidean antara dua posisi
            distance = math.hypot(current_x - prev_x, current_y - prev_y)
            
            # Menghitung jarak nyata (dalam meter) berdasarkan skala
            distance_in_meters = distance * scale_factor  # Menggunakan scale_factor untuk konversi ke meter
            
            # Kecepatan dalam km/h
            fps = cap.get(cv2.CAP_PROP_FPS)
            speed_in_kmh = distance_in_meters * fps * 3.6  # Konversi ke km/h

            return speed_in_kmh
        else:
            return 0

    def update_previous_positions(self, id, current_position):
        """Memperbarui posisi kendaraan untuk frame selanjutnya"""
        self.previous_positions[id] = current_position

tracker = EuclideanDistTracker()

while True:
    ret, frame = cap.read()
    if not ret:
        break

    # Jalankan inferensi untuk mendeteksi objek dengan YOLO
    results = model(frame)

    # Ambil bounding boxes dan labels dari hasil deteksi
    boxes = results[0].boxes.xyxy.numpy()  # Ambil bounding boxes (xyxy)
    confidences = results[0].boxes.conf.numpy()  # Ambil confidence
    classes = results[0].boxes.cls.numpy()  # Ambil class indices

    # Buat deteksi untuk tracker (x1, y1, x2, y2, id)
    detections = []
    for box, conf, cls in zip(boxes, confidences, classes):
        class_name = model.names[int(cls)]  # Ambil nama kelas
        if class_name in ['car', 'motorcycle']:  # Fokus pada kendaraan saja
            x1, y1, x2, y2 = box
            detection = [int(x1), int(y1), int(x2), int(y2)]
            detections.append(detection)

    # Update tracker dengan deteksi baru
    objects_bbs_ids = tracker.update(detections)

    # Proses objek yang terdeteksi dan hitung kecepatannya
    for obj in objects_bbs_ids:
        x, y, w, h, id = obj
        current_position = (x + w // 2, y + h // 2)  # Titik pusat objek

        # Hitung kecepatan kendaraan
        sp = tracker.getsp(id, current_position)

        # Tampilkan ID dan kecepatan kendaraan di atas bounding box
        cv2.putText(frame, f"#{id} {sp:.2f} km/h", (x, y - 10), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 2)

        # Jika kendaraan melebihi batas kecepatan, beri warna merah pada bounding box
        if sp > limit:
            cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 0, 255), 2)  # Merah
        else:
            cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)  # Hijau

        # Update posisi kendaraan untuk frame berikutnya
        tracker.update_previous_positions(id, current_position)

    # Tampilkan frame dengan deteksi dan pelacakan
    cv2.imshow("Frame", frame)

    # Keluar dari loop jika menekan ESC
    if cv2.waitKey(1) & 0xFF == 27:
        break

cap.release()
cv2.destroyAllWindows()
