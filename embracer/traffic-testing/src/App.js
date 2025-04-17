import './App.css';
import React, { useState } from 'react';


function App() {
  const [history, setHistory] = useState([Array(9).fill(null)]);
  const [stepNumber, setStepNumber] = useState(0);
  const [nextPlayer, setNextPlayer] = useState('X');
  const [winner, setWinner] = useState(null);

  const handleClick = (index) => {
    // 승자가 있거나 클릭된 칸에 값이 있다면 무시
  if (winner || history[stepNumber][index]) return;

  const currentBoard = history[stepNumber];
  const newBoard = [...currentBoard];
  newBoard[index] = nextPlayer;

  const newHistory = history.slice(0, stepNumber + 1); // 타임 트래블 후 새로 둔 경우 기존 이후 삭제
  newHistory.push(newBoard);

  setHistory(newHistory);
  setStepNumber(newHistory.length - 1);

  const newWinner = calculateWinner(newBoard);
  if (newWinner) {
    setWinner(newWinner);
  } else {
    setNextPlayer(prev => (prev === 'X' ? 'O' : 'X'));
  }

  };

  const resetBoard = () => {
    setHistory([Array(9).fill(null)]);
    setStepNumber(0);
    setNextPlayer('X');
    setWinner(null);
  };

  const jumpTo = (step) => {
    setStepNumber(step);
    setWinner(calculateWinner(history[step]));
    setNextPlayer(step % 2 === 0 ? 'X' : 'O');
  };

  const calculateWinner = (squares) => {
    const lines = [
      [0, 1, 2],
      [3, 4, 5],
      [6, 7, 8], // rows
      [0, 3, 6],
      [1, 4, 7],
      [2, 5, 8], // cols
      [0, 4, 8],
      [2, 4, 6]  // diagonals
    ];

    for (let [a, b, c] of lines) {
      if (squares[a] && squares[a] === squares[b] && squares[a] === squares[c]) {
        return squares[a]; // return 'X' or 'O'
      }
    }
    return null;
  };

  return (
    <div className="App">
      <header className="App-header">
        <div className="player-board">
          <div className="player-container">
            <span className="title">
              {winner ? `Winner ${winner}` : `Next Player: ${nextPlayer}`}
            </span>


          </div>
          <div className="grid">
            {history[stepNumber].map((value, index) => (
              <div
                key={index}
                className={`cell ${winner ? 'disabled' : ''}`}
                onClick={() => handleClick(index)}
              >
                {value}
              </div>
              ))}
            </div>
        </div>
        <div className="action-container">
          <div>
          <span className="player-label">1.</span> <button onClick={resetBoard}>
             Go to game start
           </button>
           </div>
           {history.map((_, move) => (
              move > 0 && (
                <div key={move}>
                  <span className="player-label">{move + 1}.</span>
                  <button onClick={() => jumpTo(move)}>
                    Go to Move #{move}
                  </button>
                </div>
              )
            ))}
        </div>
      </header>
    </div>
  );
}

export default App;
