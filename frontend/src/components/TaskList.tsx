import React from 'react';
import { Text, TouchableOpacity, View } from 'react-native';
import { Task } from '../services/api';

export function TaskList({ tasks, onEdit }: { tasks: Task[]; onEdit: (task: Task) => void }) {
  return <View>{tasks.map(task => <TouchableOpacity key={task.id} onPress={() => onEdit(task)} style={{ padding: 14, borderBottomWidth: 1, borderBottomColor: '#eee' }}><Text style={{ fontWeight: '700' }}>{task.title}</Text><Text>{task.status}</Text></TouchableOpacity>)}</View>;
}
