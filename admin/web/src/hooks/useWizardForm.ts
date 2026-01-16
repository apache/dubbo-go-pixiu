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

import { useState, useCallback } from 'react';

export interface UseWizardFormOptions<T> {
  initialValues: T;
  totalSteps: number;
  onFinish?: (values: T) => Promise<void>;
}

export interface UseWizardFormReturn<T> {
  // Step navigation
  currentStep: number;
  setCurrentStep: (step: number) => void;
  isFirstStep: boolean;
  isLastStep: boolean;
  next: () => void;
  prev: () => void;
  goToStep: (step: number) => void;

  // Form data
  formData: T;
  updateFormData: (updates: Partial<T>) => void;
  setFormData: (data: T) => void;
  reset: () => void;

  // Submission
  submitting: boolean;
  submit: () => Promise<void>;
}

export function useWizardForm<T extends object>(
  options: UseWizardFormOptions<T>
): UseWizardFormReturn<T> {
  const { initialValues, totalSteps, onFinish } = options;

  const [currentStep, setCurrentStep] = useState(0);
  const [formData, setFormData] = useState<T>(initialValues);
  const [submitting, setSubmitting] = useState(false);

  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === totalSteps - 1;

  const next = useCallback(() => {
    if (!isLastStep) {
      setCurrentStep((prev) => prev + 1);
    }
  }, [isLastStep]);

  const prev = useCallback(() => {
    if (!isFirstStep) {
      setCurrentStep((prev) => prev - 1);
    }
  }, [isFirstStep]);

  const goToStep = useCallback(
    (step: number) => {
      if (step >= 0 && step < totalSteps) {
        setCurrentStep(step);
      }
    },
    [totalSteps]
  );

  const updateFormData = useCallback((updates: Partial<T>) => {
    setFormData((prev) => ({ ...prev, ...updates }));
  }, []);

  const reset = useCallback(() => {
    setCurrentStep(0);
    setFormData(initialValues);
    setSubmitting(false);
  }, [initialValues]);

  const submit = useCallback(async () => {
    if (!onFinish) return;

    setSubmitting(true);
    try {
      await onFinish(formData);
    } finally {
      setSubmitting(false);
    }
  }, [onFinish, formData]);

  return {
    currentStep,
    setCurrentStep,
    isFirstStep,
    isLastStep,
    next,
    prev,
    goToStep,
    formData,
    updateFormData,
    setFormData,
    reset,
    submitting,
    submit,
  };
}
