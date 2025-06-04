import React, { useState, useEffect } from 'react';
import Sidebar from './layout/Sidebar';
import Navbar from './layout/Navbar';
import './Repost.css';
import EditModal from './editdata';  // Import EditModal here
import axios from 'axios';
import { jwtDecode } from 'jwt-decode';

function Repots() {
  const [search, setSearch] = useState('');
  const [data, setData] = useState([]);
  const [user, setUser] = useState(null);
  const [editId, setEditId] = useState(null);

  // Paginasi state
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [limit, setLimit] = useState(5);  // 5 data per halaman

  // Fungsi untuk mengambil data dari API
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
        params: {
          search,  // Pencarian berdasarkan filter
          limit,   // Jumlah data per halaman (5 data per halaman)
          page,    // Halaman yang sedang diminta
        },
      });

      if (response.data && Array.isArray(response.data.data)) {
        setData(response.data.data);
        setTotalCount(response.data.meta.total_records);  // Total data yang ada
        setTotalPages(response.data.meta.total_pages);    // Total halaman
      } else {
        console.error('Format data tidak valid:', response.data);
      }
    } catch (error) {
      console.error('Gagal mengambil data:', error);
    }
  };

  useEffect(() => {
    fetchData();
  }, [search, page, limit]);

  // Fungsi untuk menghapus data
  const handleDelete = async (id) => {
    if (!window.confirm('Apakah kamu yakin ingin menghapus data ini?')) return;

    try {
      const token = localStorage.getItem('token');
      if (!token) throw new Error('Token tidak ditemukan.');

      await axios.delete(`http://localhost:8001/api/data/delete/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      alert('Data berhasil dihapus.');
      fetchData();  // Setelah menghapus, ambil data lagi
    } catch (error) {
      console.error('Gagal menghapus data:', error);
      alert('Gagal menghapus data.');
    }
  };

  // Filter data berdasarkan pencarian
  const filteredData = data.filter((item) =>
    item.jeniskendaraan?.toLowerCase().includes(search.toLowerCase())
  );

  // Log filtered data untuk debugging
  console.log('Filtered Data:', filteredData);

  // Cek apakah hasil pencarian lebih sedikit dari limit dan reset page jika perlu
  useEffect(() => {
    if (filteredData.length < limit) {
      // Reset to page 1 only if the current page exceeds the total pages after filtering
      if (page > totalPages) {
        setPage(1);  // Reset ke halaman pertama jika halaman sekarang lebih besar dari total halaman
      }
    }
  }, [filteredData, limit, totalPages, page]);

  // Fungsi untuk pindah ke halaman berikutnya
  const handleNextPage = () => {
    if (page < totalPages) {
      setPage(page + 1); // Pindah ke halaman berikutnya
    }
  };

  // Fungsi untuk pindah ke halaman sebelumnya
  const handlePrevPage = () => {
    if (page > 1) {
      setPage(page - 1); // Pindah ke halaman sebelumnya
    }
  };

  // Fungsi untuk pindah ke halaman tertentu
  const handlePageChange = (pageNumber) => {
    setPage(pageNumber);
  };

  const handleLimitChange = (event) => {
    setLimit(Number(event.target.value));
    setPage(1); // Reset ke halaman pertama saat limit berubah
  };

  return (
    <div className="dashboard-repots">
      <Sidebar />
      <div className="dashboard-main">
        <Navbar />
        <div className="dashboard-content">
          <h2>Daftar Pelanggaran</h2>

          <div className="top-bar">
            <div className="search-ok">
              <input
                type="text"
                placeholder="Cari berdasarkan jenis kendaraan..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="search-input"
              />
            </div>

            <div className="page-info">
              <span className='1'>
                Showing
              </span>
              <select
                value={limit}
                onChange={handleLimitChange}
                className="limit-select"
              >
                <option value={5}>5</option>
                <option value={10}>10</option>
                <option value={15}>15</option>
                <option value={20}>20</option>
                <option value={25}>25</option>
              </select>
              <span className='1'>
                Showing
              </span>
            </div>
          </div>

          <div className="table-wrapper">
            <table className="pelanggaran-table">
              <thead>
                <tr>
                  <th>No</th>
                  <th>Jenis Kendaraan</th>
                  <th>Jenis Pelanggaran</th>
                  <th>Lokasi</th>
                  <th>Kecepatan</th>
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
                      {/* Menambahkan formula untuk menghitung nomor urut secara benar */}
                      <td>{(page - 1) * limit + index + 1}</td> {/* Nomor urut sesuai halaman */}
                      <td>{item.jeniskendaraan}</td>
                      <td>{item.jenispelanggaran}</td>
                      <td>{item.lokasi}</td>
                       <td>{item.Kecepatan}</td>
                      <td>{new Date(item.date).toLocaleString()}</td>
                      <td>
                        {item.gambar ? (
                          <img
                            src={`http://localhost:8001${item.gambar.replace('/public', '')}`}
                            alt="Bukti"
                            width="80"
                            onError={(e) => { e.target.src = '/default.jpg'; }}
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

          {/* Pagination and Data Info */}
          <div className="pagination-wrapper">
            <div className="data-info">
              <span>
                Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, totalCount)} of {totalCount} entries
              </span>
            </div>

            {/* Pagination controls */}
            <div className="pagination">
              <button onClick={handlePrevPage} disabled={page === 1}>
                Previous
              </button>

              <div className="page-numbers">
                {[...Array(totalPages).keys()].map((number) => (
                  <button
                    key={number + 1}
                    onClick={() => handlePageChange(number + 1)}
                    className={page === number + 1 ? "active" : ""}
                  >
                    {number + 1}
                  </button>
                ))}
              </div>

              <button onClick={handleNextPage} disabled={page === totalPages}>
                Next
              </button>
            </div>
          </div>
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
