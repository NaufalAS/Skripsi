import React, { useState, useEffect } from 'react';
import Sidebar from './layout/Sidebar';
import Navbar from './layout/Navbar';
import './profile.css';
import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

function Profile() {
  const [userData, setUserData] = useState({
    name: '',
    email: '',
    no_telepon: '',
    alamat: '',
    foto: '/default.jpg',
  });

  const [tempData, setTempData] = useState({
    name: '',
    email: '',
    no_telepon: '',
    alamat: '',
  });

  const [selectedPhoto, setSelectedPhoto] = useState(null); // Store the selected photo
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUserProfile = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) return;

        const decodedToken = jwtDecode(token);
        const userId = decodedToken.user_id;

        const response = await axios.get(`http://localhost:8001/api/user/${userId}`, {
          headers: { Authorization: `Bearer ${token}` },
        });

        const data = response.data.data;
        setUserData({
          name: data.name,
          email: data.email,
          no_telepon: data.no_telepon,
          alamat: data.alamat,
          foto: `http://localhost:8001${data.foto.replace('/public', '')}`,
        });

        setTempData({
          name: data.name,
          email: data.email,
          no_telepon: data.no_telepon,
          alamat: data.alamat,
        });

        setLoading(false);
      } catch (error) {
        console.error('Gagal mengambil data profil:', error);
        setLoading(false);
      }
    };

    fetchUserProfile();
  }, []);

  const handleSave = async () => {
    try {
      const token = localStorage.getItem('token');
      const decodedToken = jwtDecode(token);
      const userId = decodedToken.user_id;

      const formData = new FormData();
      formData.append('fullname', tempData.name);
      formData.append('email', tempData.email);
      formData.append('phone_number', tempData.no_telepon);
      formData.append('alamat', tempData.alamat);

      if (selectedPhoto) {
        formData.append('potoprofile', selectedPhoto);
      }

      await axios.put(`http://localhost:8001/api/user/update/${userId}`, formData, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      // Update local state
      const updatedFotoUrl = selectedPhoto ? URL.createObjectURL(selectedPhoto) : userData.foto;
      setUserData(prev => ({
        ...prev,
        ...tempData,
        foto: updatedFotoUrl,
      }));

      // Kirim event ke Navbar
      window.dispatchEvent(new CustomEvent('userUpdated', {
        detail: {
          name: tempData.name,
          foto: updatedFotoUrl,
        }
      }));

      setIsEditing(false);
      setSelectedPhoto(null);
    } catch (error) {
      console.error('Gagal menyimpan perubahan:', error);
    }
  };

  const handleCancel = () => {
    setTempData({
      name: userData.name,
      email: userData.email,
      no_telepon: userData.no_telepon,
      alamat: userData.alamat,
    });
    setSelectedPhoto(null);
    setIsEditing(false);
  };

  // Handle image selection and preview
  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setSelectedPhoto(file);
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div className="dashboard-container">
      <Sidebar />
      <div className="main-content">
        <Navbar />
        <div className="profile-page">
          <div className={`profileone ${isEditing ? 'edit-profile' : ''}`}>
            <h1>{isEditing ? 'Edit Profile' : 'Profile'}</h1>
          </div>
          <div className="profile-layout">
            {/* Photo */}
            <div className="profile-photo-section">
              <div className="profile-photo-container">
                {/* Display selected photo or user photo */}
                <img
                  src={selectedPhoto ? URL.createObjectURL(selectedPhoto) : userData.foto}
                  alt="Profile"
                  className="profile-photo"
                />
                {isEditing && (
                  <div className="input-file-container">
                    {/* Hidden input file */}
                    <input
                      type="file"
                      accept="image/*"
                      onChange={handleImageChange}
                      id="profilePhotoInput"
                      style={{ display: 'none' }}
                    />
                    <button
                      onClick={() => document.getElementById('profilePhotoInput').click()}
                    >
                      Change Photo
                    </button>
                  </div>
                )}
              </div>
            </div>

            {/* Info */}
            <div className="profile-info-section">
              <div className="profile-info">
                <div className={`info-row ${isEditing ? 'editing' : ''}`}>
                  <strong>Name:</strong>
                  {isEditing ? (
                    <input
                      type="text"
                      value={tempData.name}
                      onChange={(e) => setTempData({ ...tempData, name: e.target.value })}
                      className="profile-input"
                    />
                  ) : (
                    <span>{userData.name}</span>
                  )}
                </div>
                <div className={`info-row ${isEditing ? 'editing' : ''}`}>
                  <strong>Email:</strong>
                  {isEditing ? (
                    <input
                      type="email"
                      value={tempData.email}
                      onChange={(e) => setTempData({ ...tempData, email: e.target.value })}
                      className="profile-input"
                    />
                  ) : (
                    <span>{userData.email}</span>
                  )}
                </div>
                <div className={`info-row ${isEditing ? 'editing' : ''}`}>
                  <strong>Phone Number:</strong>
                  {isEditing ? (
                    <input
                      type="text"
                      value={tempData.no_telepon}
                      onChange={(e) => setTempData({ ...tempData, no_telepon: e.target.value })}
                      className="profile-input"
                    />
                  ) : (
                    <span>{userData.no_telepon}</span>
                  )}
                </div>
                <div className={`info-row ${isEditing ? 'editing' : ''}`}>
                  <strong>Address:</strong>
                  {isEditing ? (
                    <textarea
                      value={tempData.alamat}
                      onChange={(e) => setTempData({ ...tempData, alamat: e.target.value })}
                      className="profile-input"
                    />
                  ) : (
                    <span>{userData.alamat}</span>
                  )}
                </div>

                {!isEditing && (
                  <button className="btn-edit-password" onClick={() => setIsEditing(true)}>
                    Edit Profile
                  </button>
                )}

                {isEditing && (
                  <div className="edit-buttons">
                    <button className="btn-save" onClick={handleSave}>Save</button>
                    <button className="btn-cancel" onClick={handleCancel}>Cancel</button>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Profile;
