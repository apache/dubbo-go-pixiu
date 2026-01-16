/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { useEffect, useRef, useState, memo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { EChartsOption } from 'echarts';

interface SafeEChartsProps {
  option: EChartsOption;
  style?: React.CSSProperties;
  className?: string;
}

/**
 * A wrapper around ReactECharts that handles React 18 StrictMode cleanup issues.
 * The echarts-for-react library has a known bug where it fails to properly clean up
 * ResizeObserver in StrictMode, causing "disconnect" errors.
 */
function SafeEChartsComponent({ option, style, className }: SafeEChartsProps) {
  const [mounted, setMounted] = useState(false);
  const chartRef = useRef<ReactECharts>(null);

  useEffect(() => {
    // Delay mounting to avoid StrictMode double-mount issues
    let isMounted = true;
    const timer = setTimeout(() => {
      if (isMounted) setMounted(true);
    }, 0);
    return () => {
      isMounted = false;
      clearTimeout(timer);
      setMounted(false);
    };
  }, []);

  if (!mounted) {
    return <div style={style} className={className} />;
  }

  return (
    <ReactECharts
      ref={chartRef}
      option={option}
      style={style}
      className={className}
      notMerge={true}
      lazyUpdate={true}
    />
  );
}

export const SafeECharts = memo(SafeEChartsComponent);
