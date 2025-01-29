CREATE TABLE teachers (
    teacher_id SERIAL PRIMARY KEY, 
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT FALSE
);

CREATE TABLE students (
    student_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    class_number VARCHAR(2) NOT NULL
);

CREATE TABLE code_words (
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_used BOOLEAN DEFAULT FALSE
);

CREATE TABLE teacher_student (
    teacher_student_id SERIAL PRIMARY KEY,
    teacher_id INT REFERENCES teachers(teacher_id) ON DELETE CASCADE,
    student_id INT REFERENCES students(student_id) ON DELETE CASCADE
);

CREATE TABLE assignments (
    assignment_id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    teacher_id INT REFERENCES teachers(teacher_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE assignment_files (
    id SERIAL PRIMARY KEY,
    assignment_id INT REFERENCES assignments(assignment_id) ON DELETE CASCADE,
    url TEXT NOT NULL
);

CREATE TABLE student_assignment (
    student_assignment_id SERIAL PRIMARY KEY,
    assignment_id INT REFERENCES assignments(assignment_id) ON DELETE CASCADE,
    student_id INT REFERENCES students(student_id) ON DELETE CASCADE,
    teacher_id INT REFERENCES teachers(teacher_id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deadline TIMESTAMP,
    status VARCHAR(50) DEFAULT 'не решено'
);

CREATE TABLE submissions (
    submission_id SERIAL PRIMARY KEY,
    assignment_id INT REFERENCES student_assignment(student_assignment_id) ON DELETE CASCADE,
    student_id INT REFERENCES students(student_id) ON DELETE CASCADE,
    submission_text TEXT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    graded BOOLEAN DEFAULT FALSE
);

CREATE TABLE submission_files (
    id SERIAL PRIMARY KEY,
    submission_id INT REFERENCES submissions(submission_id) ON DELETE CASCADE,
    url TEXT NOT NULL
);

CREATE TABLE grades (
    grade_id SERIAL PRIMARY KEY,
    student_id INT REFERENCES students(student_id) ON DELETE CASCADE,
    submission_id INT REFERENCES submissions(submission_id) ON DELETE CASCADE,
    grade INT CHECK (grade >= 1 AND grade <= 5),
    feedback TEXT
);