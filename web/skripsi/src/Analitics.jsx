import React from 'react';
import { Line, Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, ArcElement, Tooltip, Legend } from 'chart.js';
import Sidebar from './layout/Sidebar';
import Navbar from './layout/Navbar';
import './Analitics.css';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Tooltip,
  Legend
);

function Analitics() {
  const lineData = {
    labels: ['2025-04-21', '2025-04-22', '2025-04-23', '2025-04-24', '2025-04-25'],
    datasets: [
      {
        label: 'Pelanggaran',
        data: [31, 40, 28, 51, 42],
        fill: true,
        backgroundColor: 'rgba(0, 227, 150, 0.3)',
        borderColor: '#00E396',
        tension: 0.4
      }
    ]
  };

  const doughnutData = {
    labels: ['Motor', 'Mobil'],
    datasets: [
      {
        data: [76, 24],
        backgroundColor: ['#FEB019', '#E0E0E0'],
        hoverBackgroundColor: ['#FEB019', '#E0E0E0']
      }
    ]
  };

  return (
    <div className="dashboard-page">
      <Sidebar />
      <div className="dashboard-main">
        <Navbar />
        <div className="dashboard-content">
          <h2 className="dashboard-title">Dashboard Statistik Pelanggaran</h2>

          <div className="summary-cards">
            <div className="summary-card">
              <h4>Jumlah Pelanggaran</h4>
              <p className="summary-value">192</p>
            </div>
            <div className="summary-card">
              <h4>Rata-rata Kecepatan</h4>
              <p className="summary-value">58 km/jam</p>
            </div>
            <div className="summary-card">
              <h4>Total Kendaraan</h4>
              <p className="summary-value">320</p>
            </div>
          </div>

          <div className="charts-grid">
            <div className="chart-box chart-doughnut">
              <h3>Volume Kendaraan</h3>
              <Doughnut data={doughnutData} />
            </div>

            <div className="chart-box">
              <h3>Grafik Tren Pelanggaran</h3>
              <Line data={lineData} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Analitics;
