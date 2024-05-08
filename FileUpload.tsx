import React, { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import axios from 'axios';

interface FileWithPreview extends File {
  preview: string;
}

const FileUpload: React.FC = () => {
  const [uploadProgressPercentage, setUploadProgressPercentage] = useState<number>(0);
  const [uploadFeedbackMessage, setUploadFeedbackMessage] = useState<string>("");
  const [uploadErrorMessage, setUploadErrorMessage] = useState<string>("");

  const handleFileDrop = useCallback((acceptedFiles: File[]) => {
    const singleFile = acceptedFiles[0]; // Assuming single file upload for simplicity
    if (!singleFile) {
      setUploadErrorMessage("No file selected!");
      return;
    }

    const formData = new FormData();
    formData.append('file', singleFile);

    axios.post(`${process.env.REACT_APP_UPLOAD_ENDPOINT}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        setUploadProgressPercentage(percentCompleted);
      },
    })
    .then(() => {
      setUploadFeedbackMessage("File uploaded successfully!");
      setUploadErrorMessage(""); // Clear any error messages
    })
    .catch(() => {
      setUploadErrorMessage("Error uploading file!");
      setUploadFeedbackMessage(""); // Clear feedback message
    });
  }, []);

  const {getRootProps, getInputProps, isDragActive} = useDropzone({onDrop: handleFileDrop});

  return (
    <div {...getRootProps()} style={{ border: '2px dashed black', padding: '20px', textAlign: 'center' }}>
      <input {...getInputProps()} />
      {
        isDragActive ?
          <p>Drop the file here ...</p> :
          <p>Drag 'n' drop a file here, or click to select a file</p>
      }
      {uploadProgressPercentage > 0 && <div>Upload Progress: {uploadProgressPercentage}%</div>}
      {uploadFeedbackMessage && <div>{uploadFeedbackMessage}</div>}
      {uploadErrorMessage && <div style={{color: 'red'}}>{uploadErrorMessage}</div>}
    </div>
  );
};

export default FileUpload;