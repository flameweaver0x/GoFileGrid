import React, { useState, useEffect } from 'react';
import axios from 'axios';

interface File {
  id: string;
  name: string;
  type?: string; // Optional: Assuming backend can also provide MIME type
}

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:5000'; // Default backend URL

const App: React.FC = () => {
  const [files, setFiles] = useState<File[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [operationFeedback, setOperationFeedback] = useState<string | null>(null);

  useEffect(() => {
    fetchFiles();
  }, []);

  const fetchFiles = async () => {
    try {
      const response = await axios.get(`${BACKEND_URL}/files`);
      setFiles(response.data);
      setError(null);
    } catch (err) {
      setError("Error fetching files. Please try again.");
      console.error("Error fetching files", err);
    }
  };

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const formData = new FormData();
    if (event.target.files?.length) {
      formData.append('file', event.target.files[0]);
      try {
        await axios.post(`${BACKEND_URL}/upload`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
        fetchFiles();
        setError(null);
        setOperationFeedback("File uploaded successfully!");
      } catch (err) {
        setError("Error uploading file. Please try again.");
        console.error("Error uploading file", err);
      }
    }
  };

  const handleFileDownload = async (fileId: string) => {
    try {
      const response = await axios.get(`${BACKEND_URL}/files/${fileId}`, {
        responseType: 'blob',
      });
      const file = files.find(file => file.id === fileId);
      
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', file?.name ?? 'download');
      document.body.appendChild(link);
      link.click();

      setTimeout(() => {
        window.URL.revokeObjectURL(url);
        link.remove();
      }, 100); // Timeout slightly increased for cleanup operation
      
      setError(null);
      setOperationFeedback("File downloaded successfully!");
    } catch (err) {
      setError("Error downloading file. Please try again.");
      console.error("Error downloading file", err);
    }
  };

  const handleFileDelete = async (fileId: string) => {
    try {
      await axios.delete(`${BACKEND_URL}/files/${fileId}`);
      setFiles(currentFiles => currentFiles.filter(file => file.id !== fileId));
      setError(null);
      setOperationFeedback("File deleted successfully!");
    } catch (err) {
      setError("Error deleting file. Please try again.");
      console.error("Error deleting file", err);
    }
  };

  return (
    <div>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {operationFeedback && <p style={{ color: 'green' }}>{operationFeedback}</p>}
      <input type="file" onChange={handleFileUpload} />
      <ul>
        {files.map(file => (
          <li key={file.id}>
            {file.name}
            {file.type && file.type.startsWith('image/') && (
              <img src={`${BACKEND_URL}/files/${file.id}`} alt={file.name} style={{ width: 100, height: 'auto' }} />
            )}
            <button onClick={() => handleFileDownload(file.id)}>Download</button>
            <button onClick={() => handleFileDelete(file.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default App;