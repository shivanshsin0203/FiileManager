"use client"

import axios from "axios";
import { useState } from "react";

export default function Home() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [auth, setAuth] = useState(false);
  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.post<{ token: string }>('http://localhost:8080/login', { email });
      console.log(response)
      const token = response.data.token;

      
      localStorage.setItem('token', token);

      setMessage('Login successful!');
      setAuth(true);
    } catch (error) {
      setMessage('Login failed. Please try again.'+ error);
    }
  };
  const checkToken = async () => {
    const token = localStorage.getItem('token');
    try {
      const response = await axios.get('http://localhost:8080/validate', {
        headers: { Authorization: token },
      });
      setMessage(response.data);
      if(response.data === 'Token is valid'){
        setAuth(true);}
    } catch (error) {
      setMessage('Token is invalid or expired' + error);
    }
  };
  return (
    <div className=" flex items-center justify-center flex-col">
          <h1>Hello Login please !</h1>
          <form onSubmit={handleLogin}>
        <input
          type="email"
          placeholder="Enter your email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <button type="submit">Login</button>
      </form>
      <button onClick={checkToken}>Check token</button>
      <p>{message}</p>
      <p>{auth ? 'Authenticated' : 'Not authenticated'}</p>
    </div>
  );
}
