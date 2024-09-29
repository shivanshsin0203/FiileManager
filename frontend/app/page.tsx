"use client"
import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Upload } from 'lucide-react';

interface LoginResponse {
  token: string;
}

export default function EnhancedLoginUpload() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [auth, setAuth] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);

  useEffect(() => {
    checkToken();
  }, []);

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    try {
      const response = await axios.post<LoginResponse>('http://localhost:8080/login', { email });
      const token = response.data.token;
      localStorage.setItem('token', token);
      setMessage('Login successful!');
      setAuth(true);
    } catch (error) {
      setMessage('Login failed. Please try again.');
    }
  };

  const checkToken = async (): Promise<void> => {
    const token = localStorage.getItem('token');
    if (!token) return;
    try {
      const response = await axios.get('http://localhost:8080/validate', {
        headers: { Authorization: token },
      });
      if (response.data === 'Token is valid') {
        setAuth(true);
      }
    } catch (error) {
      setMessage('Token is invalid or expired');
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleUpload = async (): Promise<void> => {
    if (!file) {
      setMessage("Please select a file first");
      return;
    }

    try {
      // Get presigned URL from your backend
      const key= file.name + Date.now();
      const response = await axios.get(`http://localhost:8080/generatePresignedURL`, {
        params: { key: key },
      });

      const presignedURL = response.data;
      console.log("Presigned URL: ", presignedURL);

      // Track upload progress
      await axios.put(presignedURL, file, {
        onUploadProgress: (progressEvent) => {
          const percentCompleted = progressEvent.total 
            ? Math.round((progressEvent.loaded * 100) / progressEvent.total) 
            : 0;
          setUploadProgress(percentCompleted);
        },
      });

      setMessage("Video uploaded successfully!");
      await axios.post('http://localhost:8080/addQueue', {Queue:"test2",Item:key});
    } catch (error) {
      setMessage("Failed to upload video: " + error);
    }
  };

  const handleLogout = (): void => {
    localStorage.removeItem('token');
    setAuth(false);
    setMessage('');
  };

  if (!auth) {
    return (
      <div className="min-h-screen bg-gradient-to-r from-blue-400 to-purple-500 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-md w-96">
          <h1 className="text-3xl font-bold mb-6 text-center text-gray-800">Welcome Back!</h1>
          <form onSubmit={handleLogin} className="space-y-4">
            <input
              type="email"
              placeholder="Enter your email"
              value={email}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
              required
              className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <button
              type="submit"
              className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 transition duration-300"
            >
              Login
            </button>
          </form>
          {message && <p className="mt-4 text-center text-red-500">{message}</p>}
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-semibold text-gray-800">Video Uploader</h1>
            </div>
            <div className="flex items-center">
              <button
                onClick={handleLogout}
                className="ml-4 px-4 py-2 rounded bg-red-500 text-white hover:bg-red-600 transition duration-300"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="border-4 border-dashed border-gray-200 rounded-lg h-96 flex flex-col items-center justify-center">
            <input
              type="file"
              onChange={handleFileChange}
              accept="video/*"
              className="hidden"
              id="file-upload"
            />
            <label
              htmlFor="file-upload"
              className="cursor-pointer flex flex-col items-center justify-center"
            >
              <Upload size={48} className="text-gray-400 mb-4" />
              <span className="text-sm text-gray-500">
                {file ? file.name : 'Click to upload a video'}
              </span>
            </label>
            {file && (
              <button
                onClick={handleUpload}
                className="mt-4 px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition duration-300"
              >
                Upload Video
              </button>
            )}
            {uploadProgress > 0 && (
              <div className="w-full max-w-xs mt-4">
                <div className="bg-gray-200 rounded-full h-2.5">
                  <div
                    className="bg-blue-600 h-2.5 rounded-full"
                    style={{ width: `${uploadProgress}%` }}
                  ></div>
                </div>
                <div className="flex justify-between mt-2 text-gray-600 text-sm">
                  <span>0%</span>
                  
                  <span>100%</span>
                </div>
                <p className="text-sm text-gray-500 mt-2 text-center">{uploadProgress}% uploaded</p>
              </div>
            )}
          </div>
          {message && (
            <p className="mt-4 text-center text-green-500 font-semibold">{message}</p>
          )}
        </div>
      </main>
    </div>
  );
}
