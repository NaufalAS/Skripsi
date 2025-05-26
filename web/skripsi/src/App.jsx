import { useState } from 'react';
import { useNavigate } from 'react-router-dom'; // Make sure to import useNavigate
import axios from 'axios';
import './App.css';

function App() {
  const [showPassword, setShowPassword] = useState(false);
  const [name, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const navigate = useNavigate(); // useNavigate to redirect to other pages

  const togglePassword = () => {
    setShowPassword(!showPassword);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!name || !password) {
      setMessage('Username dan password harus diisi!');
      return;
    }

    try {
      const response = await axios.post('http://localhost:8001/api/user/login', {
        name,
        password,
      });

      const { code, message, data } = response.data;

      if (code === 200) {
        setMessage('Login berhasil!');
        localStorage.setItem('token', data.token);
        navigate('/dashboard'); // Redirect to dashboard on successful login
      } else {
        setMessage(message || 'Login gagal.');
      }
    } catch (error) {
      setMessage(error.response?.data?.message || 'Terjadi kesalahan saat login.');
      console.error(error);
    }
  };

  return (
    <div className="container">
      <div className="left-panel">
        <h1>Selamat Datang</h1>
        <img src="public/skripsi.jpeg" alt="Logo" className="logo" />
        <p className="description">
          Sistem Pemantau Kecepatan Kendaraan untuk Keamanan & Ketertiban Lalu Lintas.
        </p>
      </div>

      <div className="right-panel">
        <div className="login-box">
          <h2 className="login-title">Masuk ke Akun Anda</h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="label">Username</label>
              <input
                type="text"
                className="input"
                placeholder="Masukkan username"
                value={name}
                onChange={(e) => setUsername(e.target.value)}
              />
            </div>

            <div className="form-group">
              <label className="label">Password</label>
              <div className="input-password-wrapper">
                <input
                  type={showPassword ? 'text' : 'password'}
                  className="input password-input"
                  placeholder="Masukkan password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                <span className="eye-icon" onClick={togglePassword}>
                  {showPassword ? '🙈' : '👁️'}
                </span>
              </div>
            </div>

            {message && (
              <p style={{ color: 'red', textAlign: 'center', marginBottom: '10px' }}>{message}</p>
            )}

            <button type="submit" className="login-button">
              Login
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

export default App;
