import React from 'react';
import { Navigate } from 'react-router-dom';

const PrivateRoute = ({ element: Element, ...rest }) => {
  const isAuthenticated = localStorage.getItem('token'); // Memeriksa apakah token ada di localStorage

  // Jika pengguna sudah login, tampilkan halaman yang diminta, jika tidak, arahkan ke halaman login
  return isAuthenticated ? <Element {...rest} /> : <Navigate to="/" />; // Arahkan ke halaman login (/) jika belum login
};

export default PrivateRoute;
