export interface TeacherAccount {
  id: string;
  name: string;
  telegramId: string;
  telegramUsername?: string;
  specialty: string;
  phone?: string;
  email?: string;
  commissionPct: number; // e.g. 70%
  assignedCourses: string[]; // Course IDs
  assignedCoursesTitles?: string[];
  totalStudents: number;
  totalLectures: number;
  totalEarnings: number;
  isActive: boolean;
  createdAt: string;
}

export interface StudentQuestion {
  id: string;
  studentId: string;
  studentName: string;
  courseId: string;
  courseTitle: string;
  lessonId: string;
  lessonTitle: string;
  timestampInVideo?: string; // e.g. "12:45"
  questionText: string;
  replyText?: string;
  repliedAt?: string;
  status: "open" | "answered";
  createdAt: string;
}

export interface TeacherCourse {
  id: string;
  title: string;
  department: string;
  semester: string;
  studentsCount: number;
  lecturesCount: number;
  price: number;
  isPublished: boolean;
}
