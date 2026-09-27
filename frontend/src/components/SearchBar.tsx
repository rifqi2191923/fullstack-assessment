import React from 'react';
import { TextInput } from 'react-native';

export function SearchBar({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  return <TextInput accessibilityLabel="Search tasks" placeholder="Search tasks..." value={value} onChangeText={onChange} style={{ borderWidth: 1, borderColor: '#ccc', borderRadius: 8, padding: 12, marginBottom: 12 }} />;
}
