package models

type Question struct {
    ID      int      `json:"id"`
    Text    string   `json:"text"`
    Options []Option `json:"options"`
}


type Option struct {
    ID        int    `json:"id"`
    Text      string `json:"text"`
    IsCorrect bool   `json:"is_correct"`
}

type Answer struct {
    QuestionID int `json:"question_id"`
    OptionID   int `json:"option_id"`
}

type Submission struct {
    UserName string   `json:"user_name"`
    Answers  []Answer `json:"answers"`
}


type Result struct {
    CorrectAnswers int `json:"correct_answers"`
    TotalQuestions  int `json:"total_questions"`
    BetterThan     int `json:"better_than"` 
}
