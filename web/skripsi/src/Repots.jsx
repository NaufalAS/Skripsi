import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { jwtDecode } from 'jwt-decode';
import Sidebar from './layout/Sidebar';
import Navbar from './layout/Navbar';
import EditModal from './editdata';
import './Repost.css';

function Repots() {
  const [search, setSearch] = useState('');
  const [data, setData] = useState([]);
  const [user, setUser] = useState(null);
  const [editId, setEditId] = useState(null);

  const fetchData = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Token tidak ditemukan. Silakan login.');

      const decoded = jwtDecode(token);
      setUser(decoded);

      const response = await axios.get('http://localhost:8001/api/data/list', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.data && Array.isArray(response.data.data)) {
        setData(response.data.data);
      } else {
        console.error('Format data tidak valid:', response.data);
      }
    } catch (error) {
      console.error('Gagal mengambil data:', error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleDelete = async (id) => {
    if (!window.confirm('Apakah kamu yakin ingin menghapus data ini?')) return;

    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Token tidak ditemukan.');

      await axios.delete(`http://localhost:8001/api/data/delete/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      alert('Data berhasil dihapus.');
      fetchData();
    } catch (error) {
      console.error('Gagal menghapus data:', error);
      alert('Gagal menghapus data.');
    }
  };

  const filteredData = data.filter((item) =>
    item.jeniskendaraan?.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="dashboard-repots">
      <Sidebar />
      <div className="dashboard-main">
        <Navbar />
        <div className="dashboard-content">
          <h2>Daftar Pelanggaran</h2>

          <div className="search-ok">
            <input
              type="text"
              placeholder="Cari berdasarkan jenis kendaraan..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="search-input"
            />
          </div>

          <table className="pelanggaran-table">
            <thead>
              <tr>
                <th>No</th>
                <th>Jenis Kendaraan</th>
                <th>Jenis Pelanggaran</th>
                <th>Lokasi</th>
                <th>Waktu</th>
                <th>Foto Bukti</th>
                <th>Aksi</th>
              </tr>
            </thead>
            <tbody>
              {filteredData.length === 0 ? (
                <tr>
                  <td colSpan="7" style={{ textAlign: 'center' }}>Tidak ada data</td>
                </tr>
              ) : (
                filteredData.map((item, index) => (
                  <tr key={item.id}>
                    <td>{index + 1}</td>
                    <td>{item.jeniskendaraan}</td>
                    <td>{item.jenispelanggaran}</td>
                    <td>{item.lokasi}</td>
                    <td>{new Date(item.date).toLocaleString()}</td>
                    <td>
                      {item.gambar ? (
                        <img
                          src={`http://localhost:8001${item.gambar.replace('/public', '')}`}
                          alt="Bukti"
                          width="80"
                          onError={(e) => {
                            e.target.src = '/default.jpg';
                          }}
                        />
                      ) : (
                        <span>Tidak ada foto</span>
                      )}
                    </td>
                    <td>
                      <button className="btn-edit" onClick={() => setEditId(item.id)}>Edit</button>
                      <button className="btn-delete" onClick={() => handleDelete(item.id)}>Hapus</button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {editId && (
        <EditModal
          id={editId}
          onClose={() => setEditId(null)}
          onUpdated={() => {
            fetchData();
            setEditId(null);
          }}
        />
      )}
    </div>
  );
}

export default Repots;
