import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8000',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

export default api;

export const signUpTeacher = (data: { name:string; email: string; password: string }) =>
  api.post('/auth/teacher/sign-up', data);

export const signInTeacher = (data: { email: string; password: string }) =>
  api.post('/auth/teacher/sign-in', data);

export const signInStudent = (data: { code: string; password: string }) =>
  api.post('/auth/student/sign-in', data);

export const signOut = () => api.post('/auth/sign-out')

export async function confirmEmail(token: string) {
  return api.get(`/auth/confirm?token=${token}`);
}

export const requestPasswordReset = (data: { email: string }) =>
  api.post('/auth/teacher/reset-password', data);

export const resetPassword = (token, data: { password: string }) =>
  api.post(`/auth/teacher/update-password?token=${token}`, data);


export const getStudents = () => api.get('/teacher/students');

export const attachStudent = (input: {
  action: string;
  name?: string;
  class_number?: number;
  code_word?: string;
}) => api.post('/teacher/students/attach', input)

export const getStudent = (id: number) => api.get(`/teacher/students/${id}`);

export const deleteStudent = (id: number) => api.delete(`/teacher/students/${id}`);

export const createAssignment = (formData: FormData) => 
  api.post('/teacher/assignments', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });

export const getAssignments = () => 
   api.get('/teacher/assignments');

export const getAssignment = (id: number) =>
  api.get(`/teacher/assignments/${id}`);

export const updateAssignment = (id: number, formData: FormData) => 
  api.put(`/teacher/assignments/${id}`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });

export const deleteAssignment = (id: number) =>
  api.delete(`/teacher/assignments/${id}`);

export const attachHomework = async (data: {
  assignment_id: number;
  student_id: number;
  title: string;
  description?: string;
  deadline: string;
}) => api.post('/teacher/homeworks/attach', data);

export const getHomeworksTeacher = (studentId?:number) => 
  api.get(studentId ? `/teacher/homeworks?id=${studentId}` : `/teacher/homeworks`);

export const getTeacherHomeworkById = (id: number) =>
  api.get(`teacher/homeworks/${id}`);

export const updateHomework = (id: number, formData: FormData) => 
  api.put(`teacher/homeworks/${id}`, formData);

export const deleteHomework = (id: number) => 
  api.delete(`teacher/homeworks/${id}`);

export const gradeHomework = (homeworkId: number, gradeData: { grade: number, feedback: string }) => 
  api.post(`/teacher/homeworks/${homeworkId}`, gradeData);

export const getTeachers = () => api.get('/student/teachers');

export const getHomeworksStudent = (teacherId?:number) => 
  api.get(teacherId ? `/student/homeworks?id=${teacherId}` : `/student/homeworks`);

export const getStudentHomeworkById = (id: number) =>
  api.get(`student/homeworks/${id}`);

export const attachStudentHomework = (id: number, formData: FormData) => 
  api.post(`/student/homeworks/${id}`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });

export const deleteStudentHomework = (id: number) => 
  api.delete(`/student/homeworks/${id}`);
;

export const updateStudentHomework = (id: number, formData: FormData) => 
  api.put(`/student/homeworks/${id}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });