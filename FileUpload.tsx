import React, { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import axios from 'axios';

interface FileWithPreview extends File {
  preview: string;
}

const FileUpload: React.FC = () => {
  const [uploadProgress, setUploadProgress] = useState<number>(0);
  const [uploadStatus, setUploadStatus] = useState<string>("");
  const [error, setError] = useState<string>("");

  const onDrop = useCallback((acceptedFiles: File[]) => {
    const file = acceptedFiles[0]; // Handle single file upload for simplicity
    if (!file) {
      setError("No file selected!");
      return;
    }

    const formData = new FormData();
    formData.append('file', file);

    axios.post(`${process.env.REACT_APP_UPLOAD_ENDPOINT}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        setUploadProgress(percentCompleted);
      },
    })
    .then((response) => {
      setUploadStatus("File uploaded successfully!");
      setError(""); // Clear any errors
    })
    .catch((error) => {
      setError("Error uploading file!");
      setUploadStatus(""); // Clear status message
    });
  }, []);

  const {getRootProps, getInputProps, isDragActive} = useDropzone({onDrop});

  return (
    <div {...getRootProps()} style={{ border: '2px dashed black', padding: '20px', textAlign: 'center' }}>
      <input {...getInputProps()} />
      {
        isDragActive ?
          <p>Drop the file here ...</p> :
          <p>Drag 'n' drop a file here, or click to select a file</p>
      }
      {uploadProgress > 0 && <div>Upload Progress: {uploadProgress}%</div>}
      {uploadStatus && <div>{uploadStatus}</div>}
      {error && <div style={{color: 'red'}}>{error}</div>}
    </div>
  );
};

export default FileUpload;