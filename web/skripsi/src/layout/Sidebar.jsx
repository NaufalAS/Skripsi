import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './Sidebar.css';

function Sidebar() {
  const location = useLocation();

  return (
    <div className="sidebar">
      <ul>
        <li>
          <Link
            to="/dashboard"
            className={location.pathname === '/dashboard' ? 'active' : ''}
          >
            <img src="performance.png" alt="Monitoring" className="menu-icon" />
            Monitoring
          </Link>
        </li>
        <li>
          <Link
            to="/analitics"
            className={location.pathname === '/analitics' ? 'active' : ''}
          >
            <img src="analitics.png" alt="Analitics" className="menu-icon" />
            Analitics
          </Link>
        </li>
        <li>
          <Link
            to="/repots"
            className={location.pathname === '/repots' ? 'active' : ''}
          >
            <img src="security-access.png" alt="Reports" className="menu-icon" />
            Reports
          </Link>
        </li>
        {/* <li>
          <Link
            to="/profile"
            className={location.pathname === '/profile' ? 'active' : ''}
          >
            <img src="setting.png" alt="Members" className="menu-icon" />
            Members
          </Link>
        </li> */}
        <li className="logout">
          <Link
            to="/profile"
            className={location.pathname === '/profile' ? 'active' : ''}
          >
            <img src="setting.png" alt="Logout" className="menu-icon" />
            Logout
          </Link>
        </li>
      </ul>
    </div>
  );
}

export default Sidebar;
