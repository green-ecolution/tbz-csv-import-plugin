import React, { useState } from "react"

function App() {
  const [file, setFile] = useState<File | null>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  }

  const handleUpload = async () => {
    if (file) {
      const formData = new FormData()
      formData.append('file', file)

      try {
        return await fetch('/api-local/v1/plugin/csv-import/upload', {
          method: 'POST',
          body: formData
        }).then(res => {
          if (res.status >= 400) {
            throw res.json()
          }
          return
        })
      } catch (error) {
        console.log(error)
      }
    }
  }

  return (
    <>
      <div>
        <input id="file" type="file" onChange={handleFileChange} />
      </div>

      <button onClick={handleUpload}>
        Upload file
      </button>
    </>
  )
}

export default App
