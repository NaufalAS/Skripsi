import cv2
import numpy as np
from collections import defaultdict
from ultralytics import YOLO
from deep_sort_realtime.deep_sort.detection import Detection
from deep_sort_realtime.deepsort_tracker import DeepSort

# Inisialisasi YOLOv8
model = YOLO("yolov8n.pt")
video_path = "Test Video.mp4"
cap = cv2.VideoCapture(video_path)

# Inisialisasi DeepSORT
deepsort = DeepSort()

track_history = defaultdict(lambda: [])

# Ambil FPS video
fps = cap.get(cv2.CAP_PROP_FPS)

# Estimasi skala: misal 50 piksel = 1 meter
meter_per_pixel = 1 / 50  # ubah sesuai kondisi sebenarnya

# Koordinat garis pertama dan kedua (garis atas dan bawah)
upper_line_y = 250
lower_line_y = 50  # Garis bawah, misalnya 50 piksel dari bawah

while cap.isOpened():
    success, frame = cap.read()
    if not success:
        break

    # Deteksi objek menggunakan YOLOv8
    results = model.track(frame, persist=True, classes=[2, 3])
    
    boxes = results[0].boxes.xywh.cpu()  # Mendapatkan bounding box dari YOLO
    scores = results[0].boxes.conf.cpu()  # Mendapatkan skor kepercayaan
    classes = results[0].boxes.cls.cpu().tolist()  # Kelas dari objek yang terdeteksi
    
    # Membuat objek Detection untuk DeepSORT
    detections = [
        Detection(box.tolist(), score.item(), class_id) 
        for box, score, class_id in zip(boxes, scores, classes)
    ]
    
    # Memperbarui pelacakan menggunakan DeepSORT
    tracks = deepsort.update_tracks(detections, frame=frame)

    # Frame salinan untuk digambar
    annotated_frame = frame.copy()

    # Filter data hasil deteksi (hanya mobil dan motor)
    filtered_data = [
        (detection.xywh, track.track_id, detection.class_id)  # Menggunakan atribut 'xywh' untuk bounding box
        for track, detection in zip(tracks, detections)
    ]

    for box, track_id, class_id in filtered_data:
        x, y, w, h = box
        x, y, w, h = int(x), int(y), int(w), int(h)
        label = "Motor" if int(class_id) == 3 else "Mobil"

        # Menggunakan DeepSORT untuk pelacakan dan menghitung kecepatan
        speed_kmph = 0
        if len(track_history[track_id]) >= 2:
            x1, y1 = track_history[track_id][-2]
            x2, y2 = track_history[track_id][-1]
            
            # Hitung jarak dalam piksel
            dist = np.sqrt((x2 - x1)**2 + (y2 - y1)**2)
            
            # Waktu antar frame dalam detik
            time_elapsed = 1 / fps  # frame per detik

            # Pastikan ada pergerakan cukup besar antar frame
            if dist > 5:  # hanya hitung kecepatan jika objek bergerak cukup jauh
                # Kecepatan dalam piksel per detik
                speed_px_per_sec = dist / time_elapsed
                
                # Mengkonversi ke kecepatan dalam km/h
                speed_kmph = speed_px_per_sec * meter_per_pixel * 3.6

        # Menambahkan ke histori kecepatan dan menghitung rata-rata
        track_history[track_id].append((x, y))
        if len(track_history[track_id]) > 5:  # Ambil rata-rata 5 frame
            track_history[track_id].pop(0)
        avg_speed = np.mean([speed_kmph for _, speed_kmph in track_history[track_id]])

        # Warna bounding box dan teks
        color_box = (0, 255, 0) if int(class_id) == 3 else (255, 0, 0)  # motor: hijau, mobil: biru
        text_color = (255, 255, 255) if int(class_id) == 3 else (0, 0, 0)

        # Gambar bounding box
        top_left = (x - w // 2, y - h // 2)
        bottom_right = (x + w // 2, y + h // 2)
        cv2.rectangle(annotated_frame, top_left, bottom_right, color_box, 2)

        # Teks kustom: ID + label + kecepatan
        custom_label = f"{label} ID: {track_id} | {int(avg_speed)} km/h"
        cv2.putText(
            annotated_frame,
            custom_label,
            (top_left[0], top_left[1] - 10),
            cv2.FONT_HERSHEY_SIMPLEX,
            0.5,
            text_color,
            2
        )

        # Jalur gerak
        points = np.array(track_history[track_id], dtype=np.int32).reshape((-1, 1, 2))
        cv2.polylines(annotated_frame, [points], isClosed=False, color=color_box, thickness=2)

    # Gambar dua garis horizontal merah
    cv2.line(annotated_frame, (0, upper_line_y), (annotated_frame.shape[1], upper_line_y), (0, 0, 255), 2)  # Garis atas
    cv2.line(annotated_frame, (0, annotated_frame.shape[0] - lower_line_y), (annotated_frame.shape[1], annotated_frame.shape[0] - lower_line_y), (0, 0, 255), 2)  # Garis bawah

    # Tampilkan frame
    cv2.imshow("Tracking Kendaraan + Kecepatan", annotated_frame)
    if cv2.waitKey(1) & 0xFF == ord("q"):
        break

cap.release()
cv2.destroyAllWindows()
