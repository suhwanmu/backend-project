import React, { createContext, useContext, useReducer } from 'react';

// 초기 상태
const initialState = {
  apiList: [],
  userCount: null,
  duration: null,
  results: [],
  summary: {
    total: 0,
    success: 0,
    minTime: null,
    maxTime: null,
  },
  loading: false,
};

// 액션 타입 정의
export const ACTIONS = {
  ADD_API: 'ADD_API',
  SET_USERS: 'SET_USERS',
  SET_DURATION: 'SET_DURATION',
  SET_LOADING: 'SET_LOADING',
  SET_RESULTS: 'SET_RESULTS',
  SET_SUMMARY: 'SET_SUMMARY',
  RESET: 'RESET',
};

// 리듀서
function trafficReducer(state, action) {
  switch (action.type) {
    case ACTIONS.ADD_API:
      return { ...state, apiList: [...state.apiList, action.payload] };
    case ACTIONS.SET_USERS:
      return { ...state, userCount: action.payload };
    case ACTIONS.SET_DURATION:
      return { ...state, duration: action.payload };
    case ACTIONS.SET_LOADING:
      return { ...state, loading: action.payload };
    case ACTIONS.SET_RESULTS:
      return { ...state, results: action.payload };
    case ACTIONS.SET_SUMMARY:
      return { ...state, summary: action.payload };
    case ACTIONS.RESET:
      return initialState;
    default:
      return state;
  }
}

// Context 생성
const TrafficContext = createContext();

export const TrafficProvider = ({ children }) => {
  const [state, dispatch] = useReducer(trafficReducer, initialState);
  return (
    <TrafficContext.Provider value={{ state, dispatch }}>
      {children}
    </TrafficContext.Provider>
  );
};

export const useTraffic = () => useContext(TrafficContext);
