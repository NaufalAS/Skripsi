import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './edit.css';

function EditModal({ id, onClose, onUpdated }) {
  const [formData, setFormData] = useState({
    jeniskendaraan: '',
    jenispelanggaran: '',
    lokasi: '',
    date: '',
    kecepatan: '',
  });

  const [imageUrl, setImageUrl] = useState('/default.jpg');
  const [imageFile, setImageFile] = useState(null);
  const fileInputRef = React.createRef();

  useEffect(() => {
    if (!id) return;
    const fetchData = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get(`http://localhost:8001/api/data/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const data = response.data.data;
        setFormData({
          jeniskendaraan: data.jeniskendaraan,
          jenispelanggaran: data.jenispelanggaran,
          lokasi: data.lokasi,
          date: data.date.slice(0, 10),
          lokasi: data.kecepatan,
        });

        if (data.gambar) {
          setImageUrl(`http://localhost:8001${data.gambar.replace('/public', '')}`);
        }
      } catch (err) {
        console.error('Gagal ambil data:', err);
      }
    };
    fetchData();
  }, [id]);

  const handleChange = (e) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    setImageFile(file);
    const reader = new FileReader();
    reader.onloadend = () => {
      setImageUrl(reader.result); // Display the image preview
    };
    if (file) {
      reader.readAsDataURL(file); // Read the file as a data URL for preview
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const formDataToSend = new FormData();
    formDataToSend.append('jeniskendaraan', formData.jeniskendaraan);
    formDataToSend.append('jenispelanggaran', formData.jenispelanggaran);
    formDataToSend.append('lokasi', formData.lokasi);
    formDataToSend.append('date', formData.date);
    formDataToSend.append('date', formData.kecepatan);

    if (imageFile) {
      formDataToSend.append('gambar', imageFile); // Append image file
    }

    try {
      const token = localStorage.getItem('token');
      await axios.put(`http://localhost:8001/api/data/update/${id}`, formDataToSend, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'multipart/form-data',
        },
      });
      alert('Data berhasil diupdate');
      onUpdated(); // trigger refresh data
      onClose(); // close the modal
    } catch (err) {
      console.error('Gagal update data:', err);
    }
  };

  const handleIconClick = () => {
    fileInputRef.current.click(); // Trigger file input click
  };

  return (
    <div className="modal-overlay-data">
      <div className="modal-content-data">
        <h2>Edit Data</h2>
        <form onSubmit={handleSubmit}>
          <div className="photo-upload-data">
            <img src={imageUrl} alt="Foto" className="profile-photo-data" />
            <img
              src="/public/update.png"
              alt="Update Foto"
              className="update-icon-data"
              onClick={handleIconClick}
            />
            <input
              type="file"
              accept="image/*"
              onChange={handleImageChange}
              ref={fileInputRef}
              style={{ display: 'none' }}
            />
          </div>
          <div className="inputnajay-data">
            <input
              type="text"
              name="jeniskendaraan"
              value={formData.jeniskendaraan}
              onChange={handleChange}
              required
            />
            <input
              type="text"
              name="jenispelanggaran"
              value={formData.jenispelanggaran}
              onChange={handleChange}
              required
            />
            <input
              type="text"
              name="lokasi"
              value={formData.lokasi}
              onChange={handleChange}
              required
            />
            <input
              type="date"
              name="date"
              value={formData.date}
              onChange={handleChange}
              required
            />
            <input
              type="text"
              name="kecepatan"
              value={formData.kecepatan}
              onChange={handleChange}
              required
            />
          </div>

          <div className="edit-buttons-data">
            <button type="submit" className="btn-save-data">
              Simpan
            </button>
            <button type="button" className="btn-cancel-data" onClick={onClose}>
              Batal
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default EditModal;
