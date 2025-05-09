// monorepo/web/prototyping/papercloud-cli/src/pages/RequestOTT.jsx
import { useState } from "react";
import { useNavigate } from "react-router";
import { authAPI } from "../services/api";

function RequestOTT() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      // Use the API service instead of direct axios call
      await authAPI.requestOTT(email);

      setSuccess(true);
      // Navigate to verify OTT page after successful request
      navigate("/verify-ott", { state: { email } });
    } catch (err) {
      console.error("Error requesting OTT:", err);
      setError(
        err.response?.data?.message ||
          err.message ||
          "Failed to request verification code",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h1>Login</h1>
      <p>Enter your email to receive a one-time verification code</p>

      {error && <p>{error}</p>}
      {success && <p>Verification code sent! Please check your email.</p>}

      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Sending..." : "Send Verification Code"}
        </button>
      </form>
    </div>
  );
}

export default RequestOTT;
