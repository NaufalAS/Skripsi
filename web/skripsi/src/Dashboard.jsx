  import React, { useEffect, useRef, useState } from 'react';
  import Sidebar from './layout/Sidebar';
  import Navbar from './layout/Navbar';
  import './Dasboard.css';

  function Dashboard() {
    const videoRef = useRef(null);
    const [result, setResult] = useState(null);
    const [workerId, setWorkerId] = useState(null);
    const [inferEngine, setInferEngine] = useState(null);

    // Memuat model menggunakan inferencejs dan publishable key
    useEffect(() => {
      const loadModel = async () => {
        try {
          if (window.inferencejs) {
            const { InferenceEngine, CVImage } = window.inferencejs;
            const inferEngineInstance = new InferenceEngine();

            // Mulai worker untuk model dan simpan workerId
            const worker = await inferEngineInstance.startWorker(
              "skripsi-oke-dsde8", // Ganti dengan nama model Anda
              "1", // Ganti dengan versi model Anda
              "rf_dXI8yFQGcBTNYy2ZSV6B3dkDoRY2" // Ganti dengan publishable key Anda
            );

            // Set workerId dan inferEngine hanya setelah worker berhasil dimulai
            setWorkerId(worker);
            setInferEngine(inferEngineInstance);
            console.log("Model has loaded and worker is ready!");
          } else {
            console.error("Inferencejs tidak terdeteksi.");
          }
        } catch (error) {
          console.error("Gagal memuat model:", error);
        }
      };

      loadModel();
    }, []); // Menjalankan sekali saat komponen pertama kali dimuat

    // Mengakses kamera dan mengatur video stream
    useEffect(() => {
      const getCamera = async () => {
        try {
          const stream = await navigator.mediaDevices.getUserMedia({ video: true });
          if (videoRef.current) {
            videoRef.current.srcObject = stream;
          }
        } catch (err) {
          console.error("Gagal mengakses kamera:", err);
        }
      };

      getCamera();

      return () => {
        if (videoRef.current) {
          const stream = videoRef.current.srcObject;
          const tracks = stream?.getTracks();
          tracks?.forEach(track => track.stop());
        }
      };
    }, []);

    // Mengirim gambar untuk inferensi menggunakan workerId
    const detectFrame = async (imageData) => {
      if (!workerId || !inferEngine) {
        console.error("workerId atau inferEngine tidak tersedia");
        return;
      }

      try {
        const inputImageElement = new Image();
        inputImageElement.src = imageData;

        // Tunggu sampai gambar selesai dimuat
        inputImageElement.onload = async () => {
          const { CVImage } = window.inferencejs;
          const inputImage = new CVImage(inputImageElement); // Gambar yang dikirimkan ke model

          // Menjalankan inferensi dengan workerId
          const predictions = await inferEngine.infer(workerId, inputImage);
          console.log("Predictions:", predictions);

          // Gambar bounding box jika ada prediksi
          if (predictions && predictions.length > 0) {
            drawBoundingBoxes(predictions);
          } else {
            // Jika tidak ada prediksi, hapus semua bounding box
            removeBoundingBoxes();
          }

          setResult(predictions); // Menyimpan hasil deteksi
        };
      } catch (error) {
        console.error("Error dalam deteksi:", error);
      }
    };

    // Fungsi untuk menggambar bounding box langsung di atas video
    const drawBoundingBoxes = (predictions) => {
      const videoElement = videoRef.current;
      const videoRect = videoElement.getBoundingClientRect(); // Mendapatkan posisi dan ukuran video
    
      // Menghapus bounding box yang sudah ada
      removeBoundingBoxes();
    
      predictions.forEach((prediction) => {
        const { bbox, class: label, confidence, color } = prediction;
    
        // Pastikan bbox memiliki koordinat yang valid
        if (bbox && bbox.x != null && bbox.y != null && bbox.width != null && bbox.height != null) {
          const { x, y, width, height } = bbox;
    
          // Menghitung skala berdasarkan ukuran video yang ditampilkan
          const scaleX = videoRect.width / videoElement.videoWidth;
          const scaleY = videoRect.height / videoElement.videoHeight;
    
          // Membuat elemen div untuk bounding box
          const box = document.createElement('div');
          box.classList.add('bounding-box');
          box.style.position = 'absolute';
    
          // Menentukan posisi yang benar, menghindari keluar batas
          const left = Math.max(0, x * scaleX);  // Jangan sampai bounding box keluar kiri
          let top = Math.max(0, y * scaleY);   // Jangan sampai bounding box keluar atas
          let right = Math.min(videoRect.width, (x + width) * scaleX); // Jangan sampai keluar kanan
          let bottom = Math.min(videoRect.height, (y + height) * scaleY); // Jangan sampai keluar bawah
    
          // Pastikan bounding box tidak keluar dari bawah video
          const finalHeight = bottom - top;
          
          // Jika bounding box terlalu besar dan keluar dari bawah, sesuaikan ukuran dan posisi
          if (finalHeight + top > videoRect.height) {
            top = videoRect.height - finalHeight;
          }
    
          // Menentukan posisi dan ukuran bounding box
          box.style.left = `${left}px`;
          box.style.top = `${top}px`; // Posisi top yang sudah disesuaikan
          box.style.width = `${right - left}px`;  // Pastikan lebar box tidak lebih besar dari video
          box.style.height = `${finalHeight}px`;  // Pastikan tinggi box tidak lebih besar dari video
          box.style.border = `3px solid ${color || '#FF0000'}`;  // Warna bounding box
          box.style.backgroundColor = 'rgba(255, 0, 0, 0.3)';  // Transparan
          box.style.zIndex = 10;  // Pastikan berada di atas video
    
          // Menambahkan label dan tingkat kepercayaan
          const labelText = document.createElement('span');
          labelText.textContent = `${label} (${(confidence * 100).toFixed(2)}%)`;
          labelText.style.position = 'absolute';
          labelText.style.top = '-20px';
          labelText.style.left = '0';
          labelText.style.color = 'white';
          labelText.style.fontSize = '14px';
          box.appendChild(labelText);
    
          // Menambahkan bounding box ke dalam kontainer video
          videoElement.parentElement.appendChild(box);
        }
      });
    };  
      //bawah percobaan
    // const drawBoundingBoxes = (predictions) => {
    //   const videoElement = videoRef.current;
    //   const videoRect = videoElement.getBoundingClientRect(); // Mengetahui posisi video di layar

    //   // Hapus bounding box sebelumnya
    //   removeBoundingBoxes();

    //   // Gambar bounding box untuk prediksi yang valid
    //   predictions.forEach((prediction) => {
    //     const { bbox, class: label, confidence, color } = prediction;

    //     // Pastikan bbox ada dan memiliki nilai valid
    //     if (bbox && bbox.x != null && bbox.y != null && bbox.width != null && bbox.height != null) {
    //       const { x, y, width, height } = bbox;

    //       // Menghitung posisi dan ukuran berdasarkan video
    //       const scaleX = videoRect.width / videoElement.videoWidth;
    //       const scaleY = videoRect.height / videoElement.videoHeight;

    //       // Membuat elemen div untuk bounding box
    //       const box = document.createElement('div');
    //       box.classList.add('bounding-box');
    //       box.style.position = 'absolute';
    //       box.style.left = `${x * scaleX}px`;
    //       box.style.top = `${y * scaleY}px`;
    //       box.style.width = `${width * scaleX}px`;
    //       box.style.height = `${height * scaleY}px`;
    //       box.style.border = `3px solid ${color || '#FF0000'}`; // Gunakan warna jika ada, default merah
    //       box.style.backgroundColor = 'rgba(255, 0, 0, 0.3)';
    //       box.style.zIndex = 10; // Pastikan box berada di atas video

    //       // Menambahkan label dan confidence
    //       const labelText = document.createElement('span');
    //       labelText.textContent = `${label} (${(confidence * 100).toFixed(2)}%)`;
    //       labelText.style.position = 'absolute';
    //       labelText.style.top = '-20px';
    //       labelText.style.left = '0';
    //       labelText.style.color = 'white';
    //       labelText.style.fontSize = '14px';
    //       box.appendChild(labelText);

    //       // Menambahkan bounding box ke video
    //       videoElement.parentElement.appendChild(box);
    //     }
    //   });
    // };

    // Fungsi untuk menghapus semua bounding box yang ada
    const removeBoundingBoxes = () => {
      const existingBoxes = document.querySelectorAll('.bounding-box');
      existingBoxes.forEach(box => box.remove());
    };

    // Mengonversi frame video ke gambar dan mengirim ke model setiap 1 detik
    const captureAndDetect = () => {
      if (videoRef.current) {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        canvas.width = videoRef.current.videoWidth;
        canvas.height = videoRef.current.videoHeight;
        ctx.drawImage(videoRef.current, 0, 0, canvas.width, canvas.height);
        const imageData = canvas.toDataURL('image/jpeg'); // Menggunakan 'image/jpeg' atau 'image/png' untuk data URL

        detectFrame(imageData); // Mengirim gambar ke model untuk deteksi
      }
    };

    // Start detection every second automatically
    useEffect(() => {
      const interval = setInterval(() => {
        if (workerId && inferEngine) {
          captureAndDetect(); // Hanya jalankan deteksi jika workerId dan inferEngine sudah tersedia
        }
      }, 1000); // Ambil gambar setiap 1 detik

      return () => clearInterval(interval);
    }, [workerId, inferEngine]); // Hanya akan dijalankan jika workerId dan inferEngine sudah di-set

    return (
      <div className="dashboard-page">
        <Sidebar />
        <div className="dashboard-main">
          <Navbar />
          <div className="dashboard-content">
            <div className="dashboard-box">
              <video ref={videoRef} autoPlay playsInline muted />
            </div>
          </div>
        </div>
      </div>
    );
  }

  export default Dashboard;
