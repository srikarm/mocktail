import { useState } from 'react';

function useToastify() {
  const [toastProps, setToastProps] = useState(defaultToastProps);
  function reset() {
    setToastProps(undefined);
  }

  function setToastPropsApiResponseHandler(response) {
    if (response?.status === 200) {
      setToastProps(defaultSuccessToast);
      return;
    }
    setToastProps(defaultErrorToast(response?.message));
  }

  function setToastPropsHandler(toastType, message) {
    setToastProps({
      toastType: toastType,
      message: message
    });
  }

  return {
    toastProps,
    reset,
    setToastPropsApiResponseHandler,
    setToastPropsHandler
  };
}
export default useToastify;

export const TOASTTYPES = {
  INFO: 'info',
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warn',
  DEFAULT: 'success'
};

const defaultToastProps = {
  toastType: TOASTTYPES.DEFAULT,
  message: 'Mocktail Rocks!'
};

export const defaultSuccessToast = {
  toastType: TOASTTYPES.SUCCESS,
  message: 'Success!'
};
export const defaultErrorToast = (message = 'Something went wrong!	') => {
  return {
    toastType: TOASTTYPES.ERROR,
    message
  };
};
