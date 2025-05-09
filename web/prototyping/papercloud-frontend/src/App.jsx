// monorepo/web/prototyping/papercloud-cli/src/App.jsx
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Link,
  Navigate,
} from "react-router";
import { AuthProvider, useAuth } from "./contexts/AuthContext";
import Register from "./pages/Register";
import RequestOTT from "./pages/RequestOTT";
import VerifyOTT from "./pages/VerifyOTT";
import CompleteLogin from "./pages/CompleteLogin";
import Home from "./pages/Home";
import Profile from "./pages/Profile"; // Import the new Profile component

// Protected route component
function ProtectedRoute({ children }) {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

// Navigation with authentication status
function Navigation() {
  const { isAuthenticated, logout } = useAuth();

  return (
    <nav>
      <ul>
        <li>
          <Link to="/">Home</Link>
        </li>
        {!isAuthenticated ? (
          <>
            <li>
              <Link to="/register">Register</Link>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
          </>
        ) : (
          <>
            <li>
              <Link to="/profile">Profile</Link>{" "}
              {/* New link to profile page */}
            </li>
            <li>
              <button onClick={logout}>Logout</button>
            </li>
          </>
        )}
      </ul>
    </nav>
  );
}

// Main App component
function AppContent() {
  const { isLoading } = useAuth();

  if (isLoading) {
    return <div>Loading authentication...</div>;
  }

  return (
    <div>
      <Navigation />

      <Routes>
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <Home />
            </ProtectedRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <ProtectedRoute>
              <Profile />
            </ProtectedRoute>
          }
        />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<RequestOTT />} />
        <Route path="/verify-ott" element={<VerifyOTT />} />
        <Route path="/complete-login" element={<CompleteLogin />} />
      </Routes>
    </div>
  );
}

// Wrap everything with the auth provider
function App() {
  return (
    <Router>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </Router>
  );
}

export default App;
