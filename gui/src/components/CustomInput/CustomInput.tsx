import React, { useState } from 'react';
import CustomImage from '@/components/ImageComponents/CustomImage';
import styles from './input.module.css';
import imagePath from '@/app/imagePath';

interface CustomInputProps {
  format: string;
  value: string;
  setter?: (value: string) => void;
  placeholder?: string;
  errorMessage?: string;
  isError?: boolean;
  cssClass?: string;
  disabled?: boolean;
  type?: 'primary' | 'secondary';
  icon?: string;
  size?: string;
  iconCSS?: string;
  alt?: string;
  endIcon?: string;
  endIconSize?: string;
  endIconBackCSS?: string;
}

const CustomInput: React.FC<CustomInputProps> = ({
  format,
  value,
  setter,
  placeholder,
  errorMessage,
  isError,
  cssClass,
  disabled,
  type = 'primary',
  icon,
  size,
  iconCSS,
  alt,
  endIcon,
  endIconSize,
  endIconBackCSS,
}) => {
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = event.target.value;
    setter(newValue);
  };

  const types = {
    primary: {
      css: styles.primary_medium,
      text: styles.primary_medium_field,
    },
    secondary: {
      css: styles.secondary_medium,
      text: styles.secondary_medium_field,
    },
  };

  const isPassword = format === 'password';
  const inputType = isPassword && showPassword ? 'text' : format;

  const togglePasswordVisibility = () => {
    setShowPassword((prevState) => !prevState);
  };

  return (
    <div id={'custom_input'} className={'flex w-full flex-col'}>
      <div
        className={`flex flex-row items-center ${
          isError ? styles.input_error : ''
        } ${types[type].css} ${cssClass}`}
      >
        {icon && (
          <CustomImage className={`${size} ${iconCSS}`} src={icon} alt={alt} />
        )}
        <input
          type={inputType}
          value={value}
          onChange={handleChange}
          placeholder={placeholder}
          className={`w-full ${types[type].text} outline-0`}
          disabled={disabled}
        />
        {isPassword ? (
          <CustomImage
            className={`${endIconSize} ${endIconBackCSS}`}
            src={
              showPassword
                ? imagePath.passwordUnhidden
                : imagePath.passwordHidden
            }
            alt={alt}
            onClick={togglePasswordVisibility}
          />
        ) : (
          endIcon && (
            <CustomImage
              className={`${endIconSize} ${endIconBackCSS}`}
              src={endIcon}
              alt={alt}
            />
          )
        )}
      </div>
      {isError && errorMessage && (
        <div className={styles.error_message}>{errorMessage}</div>
      )}
    </div>
  );
};

export default CustomInput;
