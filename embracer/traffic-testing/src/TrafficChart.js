// src/components/TrafficChart.js
import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

const TrafficChart = ({ results }) => {
  const colors = ['#8884d8', '#82ca9d', '#ffc658', '#ff8042', '#8dd1e1'];

  // 사용자 응답 시간 총합 기준 상위 5명 추출
  const sorted = [...results].sort((a, b) => {
    const aTime = a.steps.reduce((acc, step) => acc + step.time, 0);
    const bTime = b.steps.reduce((acc, step) => acc + step.time, 0);
    return bTime - aTime;
  });

  const topUsers = sorted.slice(0, 5);
  const chartData = [];

  topUsers.forEach(user => {
    let accumulated = 0;

    user.steps.forEach(step => {
      accumulated += step.time;
      const endTime = Number((accumulated / 1000).toFixed(2));
      const startTime = Number(((accumulated - step.time) / 1000).toFixed(2));

      chartData.push({
        time: startTime,
        [`user-${user.userId}`]: step.time.toFixed(2),
      });

      chartData.push({
        time: endTime,
        [`user-${user.userId}`]: step.time.toFixed(2),
      });
    });
  });

  // 평균 라인 만들기
  const timeMap = new Map();
  chartData.forEach(entry => {
    Object.entries(entry).forEach(([key, val]) => {
      if (key === 'time') return;
      const timeKey = entry.time;
      if (!timeMap.has(timeKey)) {
        timeMap.set(timeKey, []);
      }
      timeMap.get(timeKey).push(Number(val));
    });
  });

  timeMap.forEach((arr, time) => {
    const avg = arr.reduce((a, b) => a + b, 0) / arr.length;
    chartData.push({ time, average: avg.toFixed(2) });
  });

  // 정렬
  const sortedChartData = chartData.sort((a, b) => a.time - b.time);

  const maxTime = Math.max(...sortedChartData.map(d => d.time));

  return (
    <ResponsiveContainer width='100%' height={300}>
      <LineChart data={sortedChartData}>
        <CartesianGrid strokeDasharray='3 3' />
        <XAxis
          dataKey='time'
          type='number'
          domain={[0, maxTime]}
          tickFormatter={tick => `${tick}s`}
          label={{
            value: '응답 완료 시간 (초)',
            position: 'insideBottomRight',
            offset: -5,
          }}
        />
        <YAxis
          label={{
            value: '응답 시간 (ms)',
            angle: -90,
            position: 'insideLeft',
          }}
        />
        <Tooltip />
        <Legend />
        {topUsers.map((user, i) => (
          <Line
            key={user.userId}
            type='linear'
            dataKey={`user-${user.userId}`}
            stroke={colors[i % colors.length]}
            dot={false}
            isAnimationActive={false}
          />
        ))}
        <Line
          type='monotone'
          dataKey='average'
          stroke='#ff7300'
          strokeWidth={2}
          dot={false}
          name='Average Response Time'
        />
      </LineChart>
    </ResponsiveContainer>
  );
};
export default TrafficChart;
