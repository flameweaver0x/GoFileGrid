import React, { useState, useEffect } from 'react';
import axios from 'axios';

interface File {
  id: string;
  name: string;
}

const App: React.FC = () => {
  const [files, setFiles] = useState<File[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchFiles();
  }, []);

  const fetchFiles = async () => {
    try {
      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/files`);
      setFiles(response.data);
      setError(null); // Clear any previous errors
    } catch (error) {
      const message = "Error fetching files. Please try again.";
      console.error(message, error);
      setError(message);
    }
  };

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const formData = new FormData();
    if (event.target.files?.length) {
      formData.append('file', event.target.files[0]);
      try {
        await axios.post(`${process.env.REACT_APP_BACKEND_URL}/upload`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data', // Corrected the content type
          },
        });
        fetchFiles();
        setError(null); // Clear any previous errors
      } catch (error) {
        const message = "Error uploading file. Please try again.";
        console.error(message, error);
        setError(message);
      }
    }
  };

  const handleFileDownload = async (fileId: string) => {
    try {
      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/files/${fileId}`, {
        responseType: 'blob',
      });
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', files.find(file => file.id === fileId)?.name ?? '');
      document.body.appendChild(link);
      link.click();
      setError(null); // Clear any previous errors
    } catch (error) {
      const message = "Error downloading file. Please try again.";
      console.error(message, error);
      setError(message);
    }
  };

  const handleFileDelete = async (fileId: string) => {
    try {
      await axios.delete(`${process.env.REACT_APP_BACKEND_URL}/files/${fileId}`);
      fetchFiles();
      setError(null); // Clear any previous errors
    } catch (error) {
      const message = "Error deleting file. Please try again.";
      console.error(message, error);
      setError(message);
    }
  };

  return (
    <div>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <input type="file" onChange={handleFileUpload} />
      <ul>
        {files.map(file => (
          <li key={file.id}>
            {file.name}
            <button onClick={() => handleFileDownload(file.id)}>Download</button>
            <button onClick={() => handleFileDelete(file.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default App;