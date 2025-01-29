import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Header from './components/Header.tsx';
import HomePage from './components/HomePage.tsx'
import SignUp from './components/Auth/SignUp.tsx';
import SignIn from './components/Auth/SignIn.tsx';
import ConfirmEmail from './components/Auth/ConfirmEmail.tsx';
import ResetPassword from './components/Auth/ResetPassword.tsx';
import UpdatePassword from './components/Auth/UpdatePassword.tsx';
import StudentsList from './components/Teacher/StudentList.tsx';
import AttachStudent from './components/Teacher/AttachStudent.tsx';
import HomeworkTeacherList from './components/Teacher/HomeworkList.tsx';
import AssignmentsList from './components/Teacher/AssignmentsList.tsx';
import CreateAssignment from './components/Teacher/CreateAssignment.tsx';
import AssignmentDetails from './components/Teacher/AssignmentDetails.tsx';
import NotFoundPage from './components/NotFoundPage.tsx';
import AttachAssignment from './components/Teacher/AttachAssignment.tsx';
import TeachersList from './components/Student/TeachersList.tsx';
import HomeworkTeacherDetails from './components/Teacher/HomeworkDetails.tsx';
import HomeworkStudentDetails from './components/Student/HomeworkDetails.tsx';
import Homeworks from './components/Student/HomeworkList.tsx';

export default function App() {
  return (
    <Router>
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/auth/teacher/sign-up" element={<SignUp />} />
          <Route path="/auth/teacher/sign-in" element={<SignIn role="teacher" />} />
          <Route path="/auth/student/sign-in" element={<SignIn role="student" />} />
          <Route path="/auth/confirm" element={<ConfirmEmail />} />
          <Route path="/auth/teacher/update-password" element={<UpdatePassword />} />
          <Route path="auth/teacher/reset-password" element={<ResetPassword />} />

          <Route path="/teacher/students" element={<StudentsList />} />
          <Route path="/teacher/students/attach" element={<AttachStudent />} />
          <Route path="/teacher/assignments" element={<AssignmentsList />} />
          <Route path="/teacher/assignments/create" element={<CreateAssignment />} />
          <Route path="/teacher/assignments/:id" element={<AssignmentDetails />} />

          <Route path="/teacher/homeworks/attach" element={<AttachAssignment />} />
          <Route path="/teacher/homeworks" element={<HomeworkTeacherList />} />
          <Route path="/teacher/homeworks/:id" element={<HomeworkTeacherDetails />} />

          <Route path="/student/teachers" element={<TeachersList />} />
          <Route path="/student/homeworks" element={<Homeworks />} />
          <Route path="/student/homeworks/:id" element={<HomeworkStudentDetails />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>
    </Router>
  );
}