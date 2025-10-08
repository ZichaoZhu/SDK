## addStudent: 
echo "Adding student..."
curl -X POST http://localhost:9099/addStudent -H "Content-Type: application/json" -d '{"school":"S1","major":"CS","id":1001,"name":"Alice"}'

## queryStudent:
echo -e "\n\nQuerying student..."
curl -X POST http://localhost:9099/queryStudent -H "Content-Type: application/json" -d '{"school":"S1","studentId":1001}'

## 审批: 
echo -e "\n\nValidating student..."
curl -X POST http://localhost:9099/validateStudent -H "Content-Type: application/json" -d '{"school":"S1","studentId":1001,"newStatus":"Approved"}'

## queryStudent:
echo -e "\n\nQuerying student again..."
curl -X POST http://localhost:9099/queryStudent -H "Content-Type: application/json" -d '{"school":"S1","studentId":1001}'

## addGrade:
echo -e "\n\nAdding grade..."
curl -X POST http://localhost:9099/addGrade -H "Content-Type: application/json" -d '{"courseName":"OS","courseId":"OS101",    "teacher":"Zhang","school":"S1","studentId":1001,"year":2024,    "score":95.5,"semester":1}'

## queryGrade:
echo -e "\n\nQuerying grade..."
curl -X POST http://localhost:9099/queryGrade -H "Content-Type: application/json" -d '{"school":"S1","studentId":1001,"courseId":"OS101","year":2024,"semester":1}'

## validateGrade:
echo -e "\n\nValidating grade..."
curl -X POST http://localhost:9099/validateGrade -H "Content-Type:
    application/json" -d '{"school":"S1","studentId":1001,"courseId":"OS101","year":2024,"semester":1,"newStatus":"Approved"}'

## queryGrade:
echo -e "\n\nQuerying grade..."
curl -X POST http://localhost:9099/queryGrade -H "Content-Type: application/json" -d '{"school":"S1","studentId":1001,"courseId":"OS101","year":2024,"semester":1}'

## addPrice
echo -e "\n\nAdding price..."
curl -X POST http://localhost:9099/addPrice \
  -H "Content-Type: application/json" \
  -d '{"school":"S1","studentId":1001,"name":"ACM奖","id":"P001","year":2024,"level":"National","institution":"ACM"}'

## validatePrice
echo -e "\n\nValidating price..."
curl -X POST http://localhost:9099/validatePrice -H "Content-Type:
    application/json" -d '{"school":"S1","studentId":1001,"priceId":"P001","year":2024,"newStatus":"Approved"}'
