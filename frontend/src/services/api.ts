export type Task = {
  id: number;
  title: string;
  description: string;
  status: 'todo' | 'in_progress' | 'done';
  assignee_id?: number;
  created_at: string;
  updated_at: string;
};

export type TaskList = {
  items: Task[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
};

const API_URL =
  process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api';

async function request<T>(
  path: string,
  init?: RequestInit
): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
    },
    ...init,
  });

  const body = await res.json();

  if (!res.ok) {
    throw new Error(body?.error?.message || 'Request failed');
  }

  return body;
}

export function getTasks(params: {
  keyword: string;
  status: string;
  page: number;
  limit: number;
}) {
  const q = new URLSearchParams({
    page: String(params.page),
    limit: String(params.limit),
    sort: 'created_at_desc',
  });

  if (params.keyword) q.set('keyword', params.keyword);
  if (params.status) q.set('status', params.status);

  return request<TaskList>(`/tasks?${q.toString()}`);
}

export function updateTask(
  id: number,
  payload: Pick<Task, 'title' | 'description' | 'status'> & {
    assignee_id?: number;
  }
) {
  return request<Task>(`/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
}

// DELETE TASK
export function deleteTask(id: number) {
  return request<{ message: string }>(`/tasks/${id}`, {
    method: 'DELETE',
  });
}