import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import App from './App.jsx';  // Halaman login
import Dashboard from './Dashboard.jsx';
import Repots from './Repots.jsx';
import Analitics from './Analitics.jsx';
import Profile from './ptofile.jsx';
import ChangePassword from './password.jsx';
import PrivateRoute from './PrivateRoute.jsx'; // Import PrivateRoute
import EditForm from './editdata.jsx';

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <Router>
      <Routes>
        {/* Halaman login, bisa menggunakan / atau path lain */}
        <Route path="/" element={<App />} />  {/* Public Route (Login Page) */}
        
        {/* Protected Routes */}
        <Route path="/dashboard" element={<PrivateRoute element={Dashboard} />} />
        <Route path="/repots" element={<PrivateRoute element={Repots} />} />
        <Route path="/analitics" element={<PrivateRoute element={Analitics} />} />
        <Route path="/profile" element={<PrivateRoute element={Profile} />} />
        <Route path="/password" element={<PrivateRoute element={ChangePassword} />} />
      </Routes>
    </Router>
  </StrictMode>
);
