import React from 'react';
import styles from './Input.module.css';

interface InputProps {
  label: string;
  name: string;
  type?: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  error?: string;
  placeholder?: string;
}

export const Input: React.FC<InputProps> = ({
                                              label, name, type = 'text', value, onChange, error, placeholder
                                            }) => (
    <div className={styles.container}>
      <label className={styles.label} htmlFor={name}>{label}</label>
      <input
          className={`${styles.input} ${error ? styles.errorBorder : ''}`}
          id={name} name={name} type={type} value={value}
          onChange={onChange} placeholder={placeholder}
      />
      {error && <span className={styles.errorText}>{error}</span>}
    </div>
);