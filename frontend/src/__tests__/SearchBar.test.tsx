import React from 'react';
import { fireEvent, render } from '@testing-library/react-native';
import { SearchBar } from '../components/SearchBar';

test('search input emits changed keyword', () => {
  const onChange = jest.fn();
  const { getByLabelText } = render(<SearchBar value="" onChange={onChange} />);
  fireEvent.changeText(getByLabelText('Search tasks'), 'website');
  expect(onChange).toHaveBeenCalledWith('website');
});
