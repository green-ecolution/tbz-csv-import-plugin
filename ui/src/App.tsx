import { useState } from "react"

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <p>It works!</p>
      <p>Count: {count}</p>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </>
  )
}

export default App
