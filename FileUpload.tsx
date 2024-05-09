import React, { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import axios from 'axios';

interface FileWithPreview extends File {
  preview: string;
  uploadProgress?: number;
  error?: boolean;
}

const FileUpload: React.FC = () => {
  const [files, setFiles] = useState<FileWithPreview[]>([]);
  
  const uploadFile = (file: FileWithPreview, index: number) => {
    const formData = new FormData();
    formData.append('file', file);

    const config = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent: { loaded: number; total: number; }) => {
        const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        setFiles((prevFiles) =>
          prevFiles.map((f, idx) => (idx === index ? { ...f, uploadProgress: percentCompleted } : f)),
        );
      },
    };

    axios.post(`${process.env.REACT_APP_UPLOAD_ENDPOINT}`, formData, config)
      .then(() => {
        setFiles((prevFiles) =>
          prevFiles.map((f, idx) => (idx === index ? { ...f, uploadProgress: 100 } : f)),
        );
      })
      .catch(() => {
        setFiles((prevFiles) =>
          prevFiles.map((f, idx) => (idx === index ? { ...f, error: true } : f)),
        );
      });
  };

  const handleFileDrop = useCallback((acceptedFiles: File[]) => {
    const mappedFiles: FileWithPreview[] = acceptedFiles.map(file => Object.assign(file, { preview: URL.createObjectURL(file), uploadProgress: 0, error: false }));
    setFiles(prevFiles => [...prevFiles, ...mappedFiles]);
    mappedFiles.forEach(uploadFile);
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({ onDrop: handleFileDrop });

  const retryUpload = (file: FileWithPreview, index: number) => {
    setFiles((prevFiles) =>
      prevFiles.map((f, idx) => (idx === index ? { ...f, error: false, uploadProgress: 0 } : f)),
    );
    uploadFile(file, index);
  };

  const removeFile = (index: number) => {
    setFiles(prevFiles => prevFiles.filter((_, idx) => idx !== index));
  };

  return (
    <div {...getRootProps()} style={{ border: '2px dashed black', padding: '20px', textAlign: 'center' }}>
      <input {...getInputProps()} />
      {
        isDragActive ?
          <p>Drop the files here ...</p> :
          <p>Drag 'n' drop files here, or click to select files</p>
      }
      <div>
        {files.map((file, index) => (
          <div key={index}>
            {file.preview && file.type.startsWith("image/") ? (
              <img src={file.preview} alt="preview" style={{ width: "100px", height: "100px" }} />
            ) : (
              <p>{file.name}</p>
            )}
            {file.uploadProgress && <div>Upload Progress: {file.uploadProgress}%</div>}
            {file.error && <div style={{color: 'red'}}>Error uploading file! <button onClick={() => retryUpload(file, index)}>Retry</button></div>}
            <button onClick={() => removeFile(index)}>Remove</button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default FileUpload;