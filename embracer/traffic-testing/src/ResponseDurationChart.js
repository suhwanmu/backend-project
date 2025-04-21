// src/components/ResponseDurationChart.js
import React from 'react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LabelList,
} from 'recharts';

const ResponseDurationChart = ({ results }) => {
  if (!results || results.length === 0) return null;

  // 사용자별 총 응답 시간 계산
  const chartData = results.map(user => {
    const totalDuration = user.steps.reduce((sum, step) => sum + step.time, 0);
    return {
      name: `user-${user.userId}`,
      duration: Number(totalDuration.toFixed(2)),
      success: user.steps.every(s => s.success),
    };
  });

  // 응답시간 기준 내림차순 정렬 후 상위 5개 선택
  const top5 = chartData.sort((a, b) => b.duration - a.duration).slice(0, 5);

  const maxDuration = Math.max(...top5.map(d => d.duration));

  return (
    <ResponsiveContainer width='100%' height={top5.length * 40}>
      <BarChart
        data={top5}
        layout='vertical'
        margin={{ top: 20, right: 30, left: 10, bottom: 20 }}
      >
        <CartesianGrid strokeDasharray='3 3' />
        <XAxis type='number' domain={[0, maxDuration]} unit='ms' />
        <YAxis
          type='category'
          dataKey='name'
          width={100}
          tick={{ fontSize: 12 }}
        />
        <Tooltip formatter={val => `${val} ms`} />
        <Bar dataKey='duration' fill='#8884d8' isAnimationActive={false}>
          <LabelList
            dataKey='duration'
            position='right'
            formatter={val => `${val}ms`}
          />
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
};

export default ResponseDurationChart;
