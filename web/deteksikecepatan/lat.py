import cv2
import numpy as np
from ultralytics import YOLO
from traker import EuclideanDistTracker  # Pastikan kamu mengimpor kelas ini

# Batas kecepatan (km/h)
limit = 80

# Inisialisasi video capture
cap = cv2.VideoCapture("20km.mp4")

# Inisialisasi model YOLO
model = YOLO("yolov8n.pt")  # Model YOLO pre-trained

# Inisialisasi tracker (memastikan bahwa EuclideanDistTracker sudah ada di kode)
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

    # Debug: Menampilkan bounding box yang dihasilkan YOLO
    for box in boxes:
        print("Bounding Box:", box)

    # Buat deteksi untuk tracker (x1, y1, x2, y2, id)
    detections = []
    for box, conf, cls in zip(boxes, confidences, classes):
        class_name = model.names[int(cls)]  # Ambil nama kelas
        if class_name in ['car', 'motorcycle']:  # Fokus pada kendaraan saja
            x1, y1, x2, y2 = box

            # Menyesuaikan skala bounding box jika terlalu besar
            scale_factor = 0.85  # Mengurangi ukuran bounding box 10% (sesuaikan faktor ini sesuai kebutuhan)
            x1, y1, x2, y2 = int(x1 * scale_factor), int(y1 * scale_factor), int(x2 * scale_factor), int(y2 * scale_factor)

            detection = [int(x1), int(y1), int(x2), int(y2)]
            detections.append(detection)

    # Update tracker dengan deteksi baru
    objects_bbs_ids = tracker.update(detections)

    # Proses objek yang terdeteksi dan hitung kecepatannya
    for obj in objects_bbs_ids:
        x, y, w, h, id = obj
        sp = tracker.getsp(id)  # Hitung kecepatan kendaraan

        # Tampilkan ID dan kecepatan kendaraan di atas bounding box
        cv2.putText(frame, f"#{id} {sp} km/h", (x, y - 10), cv2.FONT_HERSHEY_SIMPLEX, 0.6, (0, 255, 0), 2)

        # Jika kendaraan melebihi batas kecepatan, beri warna merah pada bounding box
        if sp > limit:
            cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 0, 255), 2)  # Merah
        else:
            cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)  # Hijau

    # Tampilkan frame dengan deteksi dan pelacakan
    cv2.imshow("Frame", frame)

    # Keluar dari loop jika menekan ESC
    if cv2.waitKey(1) & 0xFF == 27:
        break

cap.release()
cv2.destroyAllWindows()
