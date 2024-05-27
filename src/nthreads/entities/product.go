package entities

import mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"

type Product struct {
	Id                  mongodb.ID         `json:"id" bson:"id"`
	ProductCategoryName string             `json:"productCategoryName" bson:"productCategoryName"`
	Name                string             `json:"name" bson:"name"`
	Price               string             `json:"price" bson:"price"`
	ImageUrl            string             `json:"imageUrl" bson:"imageUrl"`
	Date                string             `json:"date" bson:"date"`
	Description         string             `json:"description" bson:"description"`
	Warranty            string             `json:"warranty,omitempty" bson:"warranty"`
	Place               string             `json:"place" bson:"place"`
	AvailableColors     []string           `json:"availableColors,omitempty" bson:"availableColors"`
	AvailableSizes      []string           `json:"availableSizes,omitempty" bson:"availableSizes"`
	Color               string             `json:"color" bson:"color"`
	Size                string             `json:"size" bson:"size"`
	ExtraImages         []string           `json:"extraImages,omitempty" bson:"extraImages"`
	AvailableOffers     []AvailableOffers  `json:"availableOffers,omitempty" bson:"availableOffers"`
	QuestionsAnswers    []QuestionsAnswers `json:"questionsAnswers,omitempty" bson:"questionsAnswers"`
	ReviewData          []ReviewData       `json:"reviewData,omitempty" bson:"reviewData"`
	Tags                []string           `json:"tags,omitempty" bson:"tags"`
}

type AvailableOffers struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ConditionData string `json:"conditionData"`
}

type QuestionsAnswers struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ReviewData struct {
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	Rating      string   `json:"rating"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}
