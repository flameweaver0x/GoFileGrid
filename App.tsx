import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';

interface File {
  id: string;
  name: string;
}

const App: React.FC = () => {
  const [files, setFiles] = useState<File[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchFiles();
    return () => {
      files.forEach(file => {
        URL.revokeObjectURL(file.id);
      });
    };
  }, [files]);

  const fetchFiles = useCallback(async () => {
    try {
      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/files`);
      setFiles(response.data);
      setError(null);
    } catch (err) {
      const message = "Error fetching files. Please try again.";
      console.error(message, err);
      setError(message);
    }
  }, []);

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const formData = new FormData();
    if (event.target.files?.length) {
      formData.append('file', event.target.files[0]);
      try {
        await axios.post(`${process.env.REACT_APP_BACKEND_URL}/upload`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
        fetchFiles();
        setError(null);
      } catch (err) {
        const message = "Error uploading file. Please try again.";
        console.error(message, err);
        setError(message);
      }
    }
  };

const handleFileDownload = async (fileId: string) => {
    try {
      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/files/${fileId}`, {
        responseType: 'blob',
      });
      const file = files.find(file => file.id === fileId);
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', file?.name ?? '');
      document.body.appendChild(link);
      link.click();
      link.remove();
      URL.revokeObjectURL(url);
      setError(null);
    } catch (err) {
      const message = "Error downloading file. Please try again.";
      console.error(message, err);
      setError(message);
    }
  };

  const handleFileDelete = async (fileId: string) => {
    try {
      await axios.delete(`${process.env.REACT_APP_BACKEND_URL}/files/${fileId}`);
      fetchFiles();
      setError(null);
    } catch (err) {
      const message = "Error deleting file. Please try again.";
      console.error(message, err);
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