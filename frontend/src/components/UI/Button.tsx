import React from 'react';
import styles from './Button.module.css';

interface ButtonProps {
  variant?: 'primary' | 'secondary';
  disabled?: boolean;
  type?: 'button' | 'submit' | 'reset';
  children: React.ReactNode;
  onClick?: (e: React.MouseEvent<HTMLButtonElement>) => void;
}

export const Button: React.FC<ButtonProps> = ({
                                                variant = 'primary', disabled = false, type = 'button', children, onClick
                                              }) => (
    <button
        className={`${styles.button} ${styles[variant]}`}
        disabled={disabled} type={type} onClick={onClick}
    >
      {children}
    </button>
);