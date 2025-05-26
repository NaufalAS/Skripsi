import React, { useState, useEffect } from 'react';
import './Navbar.css';
import { jwtDecode } from 'jwt-decode';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

function Navbar() {
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [userName, setUserName] = useState('');
  const [profilePhoto, setProfilePhoto] = useState('/default.jpg');
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const toggleDropdown = () => setDropdownOpen(!dropdownOpen);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/');
  };

  useEffect(() => {
    const fetchUserProfile = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          setLoading(false);
          return;
        }
        const decodedToken = jwtDecode(token);
        const userId = decodedToken.user_id;

        const response = await axios.get(`http://localhost:8001/api/user/${userId}`, {
          headers: { Authorization: `Bearer ${token}` },
        });

        const userData = response.data.data;
        setUserName(userData.name);
        setProfilePhoto(`http://localhost:8001${userData.foto.replace('/public', '')}`);

        setLoading(false);
      } catch (error) {
        console.error('Error fetching profile:', error);
        setLoading(false);
      }
    };

    fetchUserProfile();

    // Pasang event listener untuk update user dari Profile
    const handleUserUpdated = (e) => {
      const { name, foto } = e.detail;
      setUserName(name);
      setProfilePhoto(foto);
    };

    window.addEventListener('userUpdated', handleUserUpdated);

    return () => {
      window.removeEventListener('userUpdated', handleUserUpdated);
    };
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div className="navbar">
      <div className="navbar-left">
        <span className="logo">📷 KameraKu</span>
      </div>
      <div className="navbar-right">
        <div className="user-info" onClick={toggleDropdown}>
          <img
            src={profilePhoto}
            alt={userName}
            className="profile-pic"
            onError={e => {
              e.target.onerror = null;
              e.target.src = '/default.jpg';
            }}
          />
          <span className="user-name-navbar">{userName}</span>
        </div>
        <div className={`dropdown ${dropdownOpen ? 'open' : ''}`}>
          <div className="dropdown-header">
            <span className="user-name">{userName}</span>
          </div>
          <ul>
            <li><a href="/profile"><i className="icon">👤</i> My Profile</a></li>
            <li><a href="/password"><i className="icon">✏️</i> Edit Password</a></li>
            <li><a onClick={handleLogout}><i className="icon">🚪</i> Log out</a></li>
          </ul>
        </div>
      </div>
    </div>
  );
}

export default Navbar;
