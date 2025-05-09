SELECT 
    users.chat_id, 
    COUNT(user_event_visits.id) AS VisitsCount
FROM 
    users 
LEFT JOIN 
    User_Event_Visits ON users.chat_id = user_event_visits.user_chat_id
GROUP BY 
    users.chat_id, users.quiz_city_name, users.date_quiz_finished
HAVING 
    COUNT(user_event_visits.id) < 4
    OR
    users.quiz_city_name is null
ORDER BY 
    users.quiz_city_name DESC, users.date_quiz_finished ASC;