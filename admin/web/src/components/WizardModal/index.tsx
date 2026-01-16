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
import { Modal, Steps, Button, message, Space } from 'antd';
import { useTranslation } from 'react-i18next';

export interface WizardStep {
  title: string;
  description?: string;
  icon?: React.ReactNode;
  content: React.ReactNode;
  validate?: () => Promise<boolean> | boolean;
}

export interface WizardModalProps {
  open: boolean;
  onClose: () => void;
  onFinish: () => Promise<void>;
  /** Callback when step changes. Receives current step index and total steps count. */
  onStepChange?: (step: number, totalSteps: number) => void;
  title: string;
  subtitle?: string;
  icon?: React.ReactNode;
  steps: WizardStep[];
  loading?: boolean;
  /** Modal width - supports number (px), string (e.g., '60vw'), or 'auto' for default */
  width?: number | string;
  /** Content height - supports number (px) or string (e.g., '50vh') */
  contentHeight?: number | string;
}

export function WizardModal({
  open,
  onClose,
  onFinish,
  onStepChange,
  title,
  steps,
  loading = false,
  width = '60vw',
  contentHeight = '55vh',
}: WizardModalProps) {
  const { t } = useTranslation();
  const [currentStep, setCurrentStep] = useState(0);
  const [stepLoading, setStepLoading] = useState(false);

  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === steps.length - 1;

  const handleNext = useCallback(async () => {
    const currentStepConfig = steps[currentStep];

    if (currentStepConfig.validate) {
      setStepLoading(true);
      try {
        const isValid = await currentStepConfig.validate();
        if (!isValid) {
          message.warning(t('common.pleaseCompleteForm'));
          setStepLoading(false);
          return;
        }
      } catch {
        message.warning(t('common.pleaseCompleteForm'));
        setStepLoading(false);
        return;
      }
      setStepLoading(false);
    }

    const nextStep = Math.min(currentStep + 1, steps.length - 1);
    setCurrentStep(nextStep);
    onStepChange?.(nextStep, steps.length);
  }, [currentStep, steps, t, onStepChange]);

  const handlePrev = useCallback(() => {
    const prevStep = Math.max(currentStep - 1, 0);
    setCurrentStep(prevStep);
    onStepChange?.(prevStep, steps.length);
  }, [currentStep, steps.length, onStepChange]);

  const handleFinish = useCallback(async () => {
    const currentStepConfig = steps[currentStep];

    if (currentStepConfig.validate) {
      setStepLoading(true);
      try {
        const isValid = await currentStepConfig.validate();
        if (!isValid) {
          message.warning(t('common.pleaseCompleteForm'));
          setStepLoading(false);
          return;
        }
      } catch {
        message.warning(t('common.pleaseCompleteForm'));
        setStepLoading(false);
        return;
      }
      setStepLoading(false);
    }

    await onFinish();
  }, [currentStep, steps, onFinish, t]);

  const handleClose = useCallback(() => {
    setCurrentStep(0);
    onStepChange?.(0, steps.length);
    onClose();
  }, [onClose, onStepChange, steps.length]);

  const stepsItems = steps.map((step, index) => ({
    key: index,
    title: step.title,
    icon: step.icon,
  }));

  return (
    <Modal
      open={open}
      onCancel={handleClose}
      title={title}
      width={width}
      destroyOnHidden
      maskClosable={false}
      footer={
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Button onClick={handleClose}>{t('common.cancel')}</Button>
          <Space>
            {!isFirstStep && (
              <Button onClick={handlePrev} disabled={stepLoading || loading}>
                {t('wizard.previous')}
              </Button>
            )}
            {isLastStep ? (
              <Button
                type="primary"
                onClick={handleFinish}
                loading={stepLoading || loading}
              >
                {t('wizard.finish')}
              </Button>
            ) : (
              <Button
                type="primary"
                onClick={handleNext}
                loading={stepLoading}
              >
                {t('wizard.next')}
              </Button>
            )}
          </Space>
        </div>
      }
    >
      <Steps
        current={currentStep}
        items={stepsItems}
        size="small"
        style={{ marginBottom: 24 }}
      />
      <div
        style={{
          height: typeof contentHeight === 'number' ? contentHeight : contentHeight,
          overflowY: 'auto',
          overflowX: 'hidden',
          paddingRight: 4,
          scrollbarWidth: 'none', // Firefox
          msOverflowStyle: 'none', // IE/Edge
        }}
        className="wizard-modal-content"
      >
        {/* Render all steps but only show the current one to preserve form state */}
        {steps.map((step, index) => (
          <div key={index} style={{ display: index === currentStep ? 'block' : 'none' }}>
            {step.content}
          </div>
        ))}
      </div>
      <style>{`
        .wizard-modal-content::-webkit-scrollbar {
          display: none;
        }
      `}</style>
    </Modal>
  );
}

export default WizardModal;
