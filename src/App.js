import { useRef, useState } from 'react';
import './App.css';
import { Page } from './page';

function App() {
  const [file, setFile] = useState(null)
  const [loading, setLoading] = useState(false)
  const [analyzedData, setAnalyzedData] = useState([])

  const fileUpload = useRef(null)
  document.title = "Passive Voice Checker"

  const customFileUpload = (e) => {
    e.preventDefault()

    fileUpload.current?.click()
  }

  const handleFileChange = (e) => {
    const selected = e.target.files?.[0];

    if (!selected) { 
      console.log("no file selected???")
      return 
    }

    setFile(selected)
  }

  const submitFile = (e) => {
    e.preventDefault()

    const formData = new FormData();
    formData.append('file', file);

    setLoading(true)

    fetch('http://localhost:8080/api/upload', {
      method: 'POST',
      body: formData,
    })
    .then(res => res.json())
    .then(data => {
      console.log("uploaded file: ", data)
      setAnalyzedData(data)
      setLoading(false)
    })
    .catch(err => {
      console.error("err")
      setLoading(false)
    })
  }

  const stopLoading = (e) => {
    e.preventDefault();
  }

  return (
    <div className="App">
        <aside>
          <h1>Passive Checker</h1>
          <br />
          
          <div className='upload-controls'>
            <div onClick={customFileUpload}  className='upload-card'> 
              <svg className='upload-icon' xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#e3e3e3"><path d="M440-320v-326L336-542l-56-58 200-200 200 200-56 58-104-104v326h-80ZM240-160q-33 0-56.5-23.5T160-240v-120h80v120h480v-120h80v120q0 33-23.5 56.5T720-160H240Z"/></svg>
                <p className='upload-text noselect'>Drop file here or click to upload</p>
                {file && <p className='upload-hint noselect'>Selected: {file.name}</p>}
            </div>
            <input ref={fileUpload} onChange={handleFileChange} type='file' accept='docx' />
            <input onClick={submitFile} className='upload-btn' type='button' value={"Upload"} />
          </div>
        </aside>
        <div className='main'>
          <div className='pages'>
            {loading && <div className="skeleton"></div>}
            {analyzedData && analyzedData.length > 0 ? (
              <Page analyzedData={analyzedData} />
            ) : (
              <p>Upload to get started!</p>
            )}
          </div>
        </div>
    </div>
  );
}

export default App;
