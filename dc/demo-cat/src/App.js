import React from 'react';
import cat1 from './assets/1.jpg';
import cat2 from './assets/2.jpg';
import cat3 from './assets/3.jpg';
import cat4 from './assets/4.jpg';
import cat5 from './assets/5.jpg';
import './App.css';

function App() {
  const [score, SetScore] = React.useState(0);
  const [pic, setPic] = React.useState(0);
  const pics = [cat1, cat2, cat3, cat4, cat5]

  function selectPic() {
    setPic(Math.floor(Math.random() * 5));
  }

  return (
    <div className="App">
        <img src={pics[pic]} alt="cat" />
        <form onSubmit={(e) => {
          e.preventDefault();
          const result = { score: parseFloat(score), pic };
          fetch('http://localhost:9090/cat', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(result)
          }).catch(alert)
          selectPic()
        }}>
        <span>Your Score:</span>
        <input type='number' step='0.01' value={score} onChange={e => SetScore(e.target.value)} />
        <input type='submit' />
        </form>
    </div>
  );
}

export default App;
