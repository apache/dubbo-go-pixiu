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

import { useState, useCallback, useMemo } from 'react';

export interface UseDynamicListReturn<T> {
  items: T[];
  setItems: (items: T[]) => void;
  add: (item: T) => void;
  addAt: (index: number, item: T) => void;
  remove: (index: number) => void;
  update: (index: number, item: T) => void;
  updatePartial: (index: number, updates: Partial<T>) => void;
  move: (fromIndex: number, toIndex: number) => void;
  clear: () => void;
  isEmpty: boolean;
  length: number;
}

export function useDynamicList<T>(initialItems: T[] = []): UseDynamicListReturn<T> {
  const [items, setItems] = useState<T[]>(initialItems);

  const add = useCallback((item: T) => {
    setItems((prev) => [...prev, item]);
  }, []);

  const addAt = useCallback((index: number, item: T) => {
    setItems((prev) => {
      const newItems = [...prev];
      newItems.splice(index, 0, item);
      return newItems;
    });
  }, []);

  const remove = useCallback((index: number) => {
    setItems((prev) => prev.filter((_, i) => i !== index));
  }, []);

  const update = useCallback((index: number, item: T) => {
    setItems((prev) => prev.map((existing, i) => (i === index ? item : existing)));
  }, []);

  const updatePartial = useCallback((index: number, updates: Partial<T>) => {
    setItems((prev) =>
      prev.map((existing, i) =>
        i === index ? { ...existing, ...updates } : existing
      )
    );
  }, []);

  const move = useCallback((fromIndex: number, toIndex: number) => {
    setItems((prev) => {
      if (
        fromIndex < 0 ||
        fromIndex >= prev.length ||
        toIndex < 0 ||
        toIndex >= prev.length
      ) {
        return prev;
      }

      const newItems = [...prev];
      const [removed] = newItems.splice(fromIndex, 1);
      newItems.splice(toIndex, 0, removed);
      return newItems;
    });
  }, []);

  const clear = useCallback(() => {
    setItems([]);
  }, []);

  const isEmpty = items.length === 0;
  const length = items.length;

  return useMemo(
    () => ({
      items,
      setItems,
      add,
      addAt,
      remove,
      update,
      updatePartial,
      move,
      clear,
      isEmpty,
      length,
    }),
    [items, add, addAt, remove, update, updatePartial, move, clear, isEmpty, length]
  );
}
