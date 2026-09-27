import React, { useCallback, useEffect, useState } from 'react';

import {
  ActivityIndicator,
  Alert,
  Button,
  Modal,
  SafeAreaView,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';

import {
  deleteTask,
  getTasks,
  Task,
  updateTask,
} from './src/services/api';

import { SearchBar } from './src/components/SearchBar';
import { TaskList } from './src/components/TaskList';

const STATUS_OPTIONS: Task['status'][] = [
  'todo',
  'in_progress',
  'done',
];

const STATUS_LABELS: Record<Task['status'], string> = {
  todo: 'To Do',
  in_progress: 'In Progress',
  done: 'Done',
};

export default function App() {
  const [keyword, setKeyword] = useState('');
  const [status, setStatus] = useState('');
  const [page, setPage] = useState(1);

  const [data, setData] = useState<{
    items: Task[];
    total_pages: number;
  }>({
    items: [],
    total_pages: 0,
  });

  const [loading, setLoading] = useState(false);
  const [editing, setEditing] = useState<Task | null>(null);
  const [saving, setSaving] = useState(false);

  // =========================
  // LOAD TASK
  // =========================
  const load = useCallback(async () => {
    setLoading(true);

    try {
      const result = await getTasks({
        keyword,
        status,
        page,
        limit: 5,
      });

      setData(result);
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : 'Gagal mengambil data task';

      Alert.alert('Error', message);
    } finally {
      setLoading(false);
    }
  }, [keyword, status, page]);

  useEffect(() => {
    load();
  }, [load]);

  // =========================
  // OPEN EDIT
  // =========================
  function openEdit(task: Task) {
    setEditing({ ...task });
  }

  // =========================
  // UPDATE TASK
  // =========================
  async function save() {
    if (!editing) {
      return;
    }

    if (!editing.title.trim()) {
      Alert.alert(
        'Validasi',
        'Judul task tidak boleh kosong.'
      );
      return;
    }

    setSaving(true);

    try {
      await updateTask(editing.id, {
        title: editing.title.trim(),
        description: editing.description || '',
        status: editing.status,
        assignee_id: editing.assignee_id,
      });

      setEditing(null);

      // Refresh list setelah update berhasil
      await load();

      Alert.alert(
        'Berhasil',
        'Task berhasil diperbarui.'
      );
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : 'Gagal memperbarui task';

      Alert.alert(
        'Gagal memperbarui task',
        message
      );
    } finally {
      setSaving(false);
    }
  }

  // =========================
  // DELETE TASK
  // =========================
  function removeTask() {
    if (!editing) {
      return;
    }

    Alert.alert(
      'Hapus Task',
      `Yakin ingin menghapus "${editing.title}"?`,
      [
        {
          text: 'Batal',
          style: 'cancel',
        },
        {
          text: 'Hapus',
          style: 'destructive',
          onPress: async () => {
            setSaving(true);

            try {
              await deleteTask(editing.id);

              setEditing(null);

              // Refresh list setelah delete berhasil
              await load();

              Alert.alert(
                'Berhasil',
                'Task berhasil dihapus.'
              );
            } catch (error) {
              const message =
                error instanceof Error
                  ? error.message
                  : 'Gagal menghapus task';

              Alert.alert(
                'Gagal menghapus task',
                message
              );
            } finally {
              setSaving(false);
            }
          },
        },
      ]
    );
  }

  return (
    <SafeAreaView
      style={{
        flex: 1,
        padding: 16,
        backgroundColor: '#f8fafc',
      }}
    >
      {/* HEADER */}
      <Text
        style={{
          fontSize: 26,
          fontWeight: '800',
          marginBottom: 16,
        }}
      >
        Tasks
      </Text>

      {/* SEARCH */}
      <SearchBar
        value={keyword}
        onChange={(value) => {
          setKeyword(value);
          setPage(1);
        }}
      />

      {/* STATUS FILTER */}
      <View
        style={{
          marginBottom: 12,
        }}
      >
        <Text
          style={{
            fontWeight: '700',
            marginBottom: 8,
          }}
        >
          Filter Status
        </Text>

        <View
          style={{
            flexDirection: 'row',
            flexWrap: 'wrap',
            gap: 8,
          }}
        >
          <Button
            title="All"
            onPress={() => {
              setStatus('');
              setPage(1);
            }}
          />

          {STATUS_OPTIONS.map((item) => (
            <Button
              key={item}
              title={STATUS_LABELS[item]}
              onPress={() => {
                setStatus(item);
                setPage(1);
              }}
            />
          ))}
        </View>
      </View>

      {/* LOADING / TASK LIST */}
      {loading ? (
        <View
          style={{
            paddingVertical: 30,
            alignItems: 'center',
          }}
        >
          <ActivityIndicator
            size="large"
            accessibilityLabel="Loading tasks"
          />

          <Text
            style={{
              marginTop: 8,
              color: '#64748b',
            }}
          >
            Memuat task...
          </Text>
        </View>
      ) : data.items.length === 0 ? (
        <View
          style={{
            paddingVertical: 40,
            alignItems: 'center',
          }}
        >
          <Text
            style={{
              color: '#64748b',
              fontSize: 16,
            }}
          >
            Tidak ada task ditemukan.
          </Text>
        </View>
      ) : (
        <TaskList
          tasks={data.items}
          onEdit={openEdit}
        />
      )}

      {/* PAGINATION */}
      <View
        style={{
          flexDirection: 'row',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginTop: 16,
        }}
      >
        <Button
          title="Previous"
          disabled={page <= 1 || loading}
          onPress={() =>
            setPage((current) => current - 1)
          }
        />

        <Text>
          Page {page} / {Math.max(data.total_pages, 1)}
        </Text>

        <Button
          title="Next"
          disabled={
            loading ||
            data.total_pages === 0 ||
            page >= data.total_pages
          }
          onPress={() =>
            setPage((current) => current + 1)
          }
        />
      </View>

      {/* EDIT MODAL */}
      <Modal
        visible={!!editing}
        transparent
        animationType="slide"
        onRequestClose={() => {
          if (!saving) {
            setEditing(null);
          }
        }}
      >
        <View
          style={{
            flex: 1,
            justifyContent: 'center',
            padding: 20,
            backgroundColor: 'rgba(0,0,0,0.35)',
          }}
        >
          <View
            style={{
              backgroundColor: '#fff',
              borderRadius: 12,
              padding: 20,
            }}
          >
            {/* MODAL TITLE */}
            <Text
              style={{
                fontSize: 20,
                fontWeight: '800',
                marginBottom: 16,
              }}
            >
              Edit Task
            </Text>

            {/* TITLE */}
            <Text
              style={{
                fontWeight: '700',
                marginBottom: 6,
              }}
            >
              Title
            </Text>

            <TextInput
              value={editing?.title || ''}
              onChangeText={(value) =>
                setEditing((current) =>
                  current
                    ? {
                        ...current,
                        title: value,
                      }
                    : current
                )
              }
              editable={!saving}
              placeholder="Title"
              style={{
                borderWidth: 1,
                borderColor: '#cbd5e1',
                borderRadius: 8,
                padding: 10,
                marginBottom: 12,
              }}
            />

            {/* DESCRIPTION */}
            <Text
              style={{
                fontWeight: '700',
                marginBottom: 6,
              }}
            >
              Description
            </Text>

            <TextInput
              value={editing?.description || ''}
              onChangeText={(value) =>
                setEditing((current) =>
                  current
                    ? {
                        ...current,
                        description: value,
                      }
                    : current
                )
              }
              editable={!saving}
              placeholder="Description"
              multiline
              numberOfLines={4}
              textAlignVertical="top"
              style={{
                borderWidth: 1,
                borderColor: '#cbd5e1',
                borderRadius: 8,
                padding: 10,
                minHeight: 90,
                marginBottom: 12,
              }}
            />

            {/* STATUS */}
            <Text
              style={{
                fontWeight: '700',
                marginBottom: 8,
              }}
            >
              Status
            </Text>

            <View
              style={{
                flexDirection: 'row',
                flexWrap: 'wrap',
                gap: 8,
                marginBottom: 20,
              }}
            >
              {STATUS_OPTIONS.map((item) => {
                const selected =
                  editing?.status === item;

                return (
                  <TouchableOpacity
                    key={item}
                    disabled={saving}
                    onPress={() =>
                      setEditing((current) =>
                        current
                          ? {
                              ...current,
                              status: item,
                            }
                          : current
                      )
                    }
                    style={{
                      paddingVertical: 9,
                      paddingHorizontal: 12,
                      borderRadius: 8,
                      borderWidth: 1,
                      borderColor: selected
                        ? '#2563eb'
                        : '#cbd5e1',
                      backgroundColor: selected
                        ? '#dbeafe'
                        : '#fff',
                    }}
                  >
                    <Text
                      style={{
                        fontWeight: selected
                          ? '700'
                          : '400',
                      }}
                    >
                      {STATUS_LABELS[item]}
                    </Text>
                  </TouchableOpacity>
                );
              })}
            </View>

            {/* MODAL BUTTONS */}
            <View
              style={{
                flexDirection: 'row',
                justifyContent: 'space-between',
                alignItems: 'center',
              }}
            >
              {/* DELETE */}
              <Button
                title="Delete"
                color="#dc2626"
                disabled={saving}
                onPress={removeTask}
              />

              {/* CANCEL + SAVE */}
              <View
                style={{
                  flexDirection: 'row',
                  gap: 12,
                }}
              >
                <Button
                  title="Cancel"
                  disabled={saving}
                  onPress={() =>
                    setEditing(null)
                  }
                />

                <Button
                  title={
                    saving
                      ? 'Saving...'
                      : 'Save'
                  }
                  disabled={saving}
                  onPress={save}
                />
              </View>
            </View>
          </View>
        </View>
      </Modal>
    </SafeAreaView>
  );
}