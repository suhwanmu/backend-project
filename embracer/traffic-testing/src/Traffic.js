import React, { useState } from 'react';
import { useTraffic, ACTIONS } from './TrafficContext';
import { simulateUser } from './LoadTester';
import { formatPercentage, formatTime } from './utils/format';
import ResponseDurationChart from './ResponseDurationChart';

function TrafficTest() {
  const { state, dispatch } = useTraffic();
  const [apiInput, setApiInput] = useState('');
  const [userInput, setUserInput] = useState('');
  const [durationInput, setDurationInput] = useState('');
  const [apiMethod, setApiMethod] = useState('GET');

  const handleApiApply = () => {
    if (apiInput.trim()) {
      dispatch({
        type: ACTIONS.ADD_API,
        payload: {
          url: apiInput.trim(),
          method: apiMethod,
        },
      });
      setApiInput('');
      setApiMethod(apiMethod);
    }
  };

  const handleUserApply = () => {
    if (userInput) {
      dispatch({ type: ACTIONS.SET_USERS, payload: Number(userInput) });
      setUserInput('');
    }
  };

  const handleDurationApply = () => {
    if (durationInput) {
      dispatch({ type: ACTIONS.SET_DURATION, payload: Number(durationInput) });
      setDurationInput('');
    }
  };

  const runLoadTest = async () => {
    const { apiList, userCount, duration } = state;
    if (apiList.length === 0 || !userCount || !duration) return;

    dispatch({ type: ACTIONS.SET_LOADING, payload: true });

    const allUsers = Array.from({ length: userCount }, (_, i) =>
      simulateUser(i, apiList)
    );
    const finalResults = await Promise.all(allUsers);

    dispatch({ type: ACTIONS.SET_RESULTS, payload: finalResults });
    dispatch({ type: ACTIONS.SET_LOADING, payload: false });

    const total = finalResults.reduce(
      (acc, user) => acc + user.steps.length,
      0
    );
    const success = finalResults.reduce(
      (acc, user) => acc + user.steps.filter(s => s.success).length,
      0
    );
    const times = finalResults.flatMap(u => u.steps.map(s => s.time));
    const minTime = times.length ? Math.min(...times) : 0;
    const maxTime = times.length ? Math.max(...times) : 0;

    dispatch({
      type: ACTIONS.SET_SUMMARY,
      payload: {
        total,
        success,
        minTime: minTime.toFixed(2),
        maxTime: maxTime.toFixed(2),
      },
    });
  };

  return (
    <div className='app'>
      <h1>Traffic Testing</h1>
      <div className='grid'>
        <div className='scenario-setup'>
          <h2>Scenario Setup</h2>
          <select
            value={apiMethod}
            onChange={e => setApiMethod(e.target.value)}
          >
            <option value='GET'>GET</option>
            <option value='POST'>POST</option>
            <option value='PUT'>PUT</option>
            <option value='DELETE'>DELETE</option>
          </select>

          <form
            onSubmit={e => {
              e.preventDefault();
              handleApiApply();
            }}
          >
            <input
              type='text'
              value={apiInput}
              onChange={e => setApiInput(e.target.value)}
              placeholder='http://test.com/api'
            />
            <button type='submit'>적용</button>
          </form>
          <div className='apply-list'>
            {state.apiList.map((api, idx) => (
              <div key={idx} className='apply-item'>
                <span>{idx + 1}.</span> [{api.method}] {api.url}
              </div>
            ))}
          </div>
        </div>

        <div className='numbers-setup'>
          <h2>Numbers Setup</h2>
          <label>
            Concurrent Users
            {state.userCount !== null && (
              <span className='red'> {state.userCount.toLocaleString()}명</span>
            )}
          </label>
          <form
            onSubmit={e => {
              e.preventDefault();
              handleUserApply();
            }}
          >
            <input
              type='number'
              value={userInput}
              onChange={e => setUserInput(e.target.value)}
              placeholder='50000'
            />
            <button onClick={handleUserApply}>적용</button>
          </form>
          <label>
            Duration (seconds)
            {state.duration !== null && (
              <span className='red'> {state.duration}초</span>
            )}
          </label>
          <form
            onSubmit={e => {
              e.preventDefault();
              handleUserApply();
            }}
          >
            <input
              type='number'
              value={durationInput}
              onChange={e => setDurationInput(e.target.value)}
              placeholder='300'
            />
            <button onClick={handleDurationApply}>적용</button>
          </form>
        </div>

        <div className='real-time-results'>
          <h2>Real-Time Results</h2>
          {state.results.length > 0 ? (
            <ResponseDurationChart results={state.results} />
          ) : (
            <div className='chart-placeholder'>Chart Placeholder</div>
          )}
        </div>

        <div className='summary'>
          <h2>Summary</h2>
          <div className='summary-item'>
            Total Requests: {state.summary.total}
          </div>
          <div className='summary-item'>
            Success Rate:{' '}
            {formatPercentage(
              state.summary.total > 0
                ? (state.summary.success / state.summary.total) * 100
                : 0
            )}
          </div>
          <div className='summary-item'>
            Min Response Time: {formatTime(state.summary.minTime)}
          </div>
          <div className='summary-item'>
            Max Response Time: {formatTime(state.summary.maxTime)}
          </div>
        </div>
      </div>
      <div className='control-buttons'>
        <button onClick={runLoadTest} disabled={state.loading}>
          Start
        </button>
        <button>Pause</button>
        <button>Stop</button>
      </div>
    </div>
  );
}

export default TrafficTest;
