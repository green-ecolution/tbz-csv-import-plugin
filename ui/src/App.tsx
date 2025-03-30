import React, { useState } from "react"

type ImportType = "create" | "update" | "archive"
const baseUrl = import.meta.env.VITE_BASE_URL ?? "api-local/v1";

export interface Tree {
  id: number
  latitude: number
  longitude: number
  number: string
  plantingYear: number
  provider: string
  readonly: boolean
  species: string
  updatedAt: string
  objectId: number
}


interface CsvTree {
  area: string
  street: string
  treeNumber: string
  species: string
  hochwert: number
  rechtswert: number
  plantingYear: number
}

interface ImportedTrees {
  tree: Tree
  importType: ImportType
}

interface ImportedTreeResponse {
  importedTrees: ImportedTrees[]
  csvTrees: CsvTree[]
}

function App() {
  const [file, setFile] = useState<File | null>(null)
  const [importedTrees, setImportedTrees] = useState<ImportedTreeResponse | null>(null)

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
        return await fetch(`/${baseUrl}/plugin/csv-import/upload`, {
          method: 'POST',
          body: formData
        }).then(res => {
          if (res.status >= 400) {
            throw res.json()
          }
          return res.json()
        }).then(data => {
          const importedTrees: ImportedTrees[] = data.importedTrees
          const csvTrees: CsvTree[] = data.rawTrees
          setImportedTrees({ importedTrees, csvTrees })
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

      <table>
        <tr>
          <th>Operation</th>
          <th>BaumNr.</th>
          <th>Gattung/Art</th>
          <th>Hochwert</th>
          <th>Rechtswert</th>
          <th>Pflanzjahr</th>
        </tr>
        {importedTrees && importedTrees.importedTrees.map(it => {
          return (
            <tr>
              <td>{it.importType}</td>
              <td>{it.tree.number}</td>
              <td>{it.tree.species}</td>
              <td>{it.tree.latitude}</td>
              <td>{it.tree.longitude}</td>
              <td>{it.tree.plantingYear}</td>
            </tr>
          )
        })}
      </table>
    </>
  )
}

export default App
