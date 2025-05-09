import { BrowserRouter as Router, Routes, Route, Link } from "react-router";
import Register from "./pages/Register";
import RequestOTT from "./pages/RequestOTT";
import VerifyOTT from "./pages/VerifyOTT";
import CompleteLogin from "./pages/CompleteLogin";
import Home from "./pages/Home";

function App() {
  return (
    <Router>
      <div>
        <nav>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/register">Register</Link>
            </li>
            <li>
              <Link to="/login">Login</Link>
            </li>
          </ul>
        </nav>

        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<RequestOTT />} />
          <Route path="/verify-ott" element={<VerifyOTT />} />
          <Route path="/complete-login" element={<CompleteLogin />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
