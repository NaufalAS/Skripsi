import React, { useState } from 'react';
import Sidebar from './layout/Sidebar';
import Navbar from './layout/Navbar';
import './password.css';
import axios from 'axios'; // Pastikan axios diinstall dan diimpor
import { jwtDecode } from 'jwt-decode'; // Import jwt-decode untuk mengambil ID dari token

function ChangePassword() {
  const [oldPassword, setOldPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [message, setMessage] = useState(''); // Untuk menampilkan pesan sukses atau error

  const handleChangePassword = async (e) => {
    e.preventDefault();

    try {
      // Ambil token dari localStorage
      const token = localStorage.getItem('token');
      if (!token) {
        setMessage("Token tidak ditemukan, harap login terlebih dahulu.");
        return;
      }

      // Ambil ID user dari token (menggunakan jwt-decode)
      const decodedToken = jwtDecode(token);
      const userId = decodedToken.user_id;

      // Siapkan data yang akan dikirim
      const payload = {
        password_lama: oldPassword,
        password_baru: newPassword
      };

      // Kirim request ke API
      const response = await axios.put(
        `http://localhost:8001/api/user/password/${userId}`,
        payload,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          }
        }
      );

      // Jika request sukses, beri pesan sukses
      if (response.status === 200) {
        setMessage("Password berhasil diperbarui.");
        alert("Password berhasil diperbarui!"); // Tambahkan alert sukses
        // Clear form setelah berhasil
        setOldPassword('');
        setNewPassword('');
      }

    } catch (error) {
      // Tangani error jika password lama salah
      if (error.response && error.response.data) {
        alert("Password lama yang Anda masukkan salah.");
      } else {
        console.error("Gagal memperbarui password:", error);
        setMessage("Gagal memperbarui password. Pastikan password lama benar.");
      }
    }
  };

  return (
    <div className="dashboard-password">
      <Sidebar />
      <div className="main-password">
        <Navbar />
        <div className="plain-center-wrapper">
          <h1 className="form-title">Ubah Password</h1>
          <p className="password-info">
            Password baru harus memiliki minimal 8 karakter, mengandung huruf besar, angka, dan simbol.
          </p>

          {/* Menampilkan pesan error atau sukses */}
          {message && <p className="message">{message}</p>}

          <form className="change-password-form" onSubmit={handleChangePassword}>
            <div className="form-group">
              <label>Password Lama</label>
              <input
                type="password"
                value={oldPassword}
                onChange={(e) => setOldPassword(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label>Password Baru</label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
              />
            </div>
            <button type="submit" className="save-btn">Simpan Password</button>
          </form>
        </div>
      </div>
    </div>
  );
}

export default ChangePassword;
